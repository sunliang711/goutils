package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL          string
	Exchanges    map[string]ExchangeConfig
	ReconnectSec int
}

type ExchangeConfig struct {
	Type   string //  "direct", "fanout", "topic" and "headers"
	Topics []string
}

type RabbitMQ struct {
	config       Config
	conn         *amqp.Connection
	ch           *amqp.Channel
	closeCh      chan struct{}
	wg           sync.WaitGroup
	reconnectMux sync.Mutex
	consumeMux   sync.Mutex

	logger *log.Logger
}

func NewRabbitMQ(config Config) (*RabbitMQ, error) {
	r := &RabbitMQ{
		config:  config,
		closeCh: make(chan struct{}),
		logger:  log.New(os.Stdout, "|RMQ| ", log.LstdFlags),
	}
	err := r.connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RabbitMQ) connect() error {
	r.reconnectMux.Lock()
	defer r.reconnectMux.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.config.URL)
	if err != nil {
		return err
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		r.conn.Close() // 确保连接关闭
		return err
	}

	for name, exchange := range r.config.Exchanges {
		r.logger.Printf("Declare exchange: %s\n", name)
		err = r.ch.ExchangeDeclare(
			name,          // name
			exchange.Type, // type
			true,          // durable
			false,         // auto-deleted
			false,         // internal
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			return err
		}
	}

	r.handleReconnect()
	return nil
}

func (r *RabbitMQ) reconnect() {
	for {
		r.logger.Println("Attempting to reconnect...")
		err := r.connect()
		if err == nil {
			r.logger.Println("Reconnected to RabbitMQ")
			return
		}
		r.logger.Printf("Failed to reconnect: %s. Retrying in %d seconds...", err, r.config.ReconnectSec)
		time.Sleep(time.Duration(r.config.ReconnectSec) * time.Second)
	}
}

func (r *RabbitMQ) handleReconnect() {
	go func() {
		for {
			err := <-r.conn.NotifyClose(make(chan *amqp.Error))
			if err != nil {
				r.logger.Printf("Connection closed: %s", err)
				r.reconnect()
			}
			// Exit the goroutine if the connection is closed
			if r.conn.IsClosed() {
				return
			}
		}
	}()
}

func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
	return r.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}
func (r *RabbitMQ) Consume(ctx context.Context, exchange, queueName string, topics []string) (<-chan amqp.Delivery, error) {
	r.consumeMux.Lock()
	defer r.consumeMux.Unlock()

	var queue amqp.Queue
	var err error
	maxRetries := 5
	retryInterval := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		r.logger.Printf("Declare queue: %s\n", queueName)
		queue, err = r.ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err == nil {
			break
		}
		r.logger.Printf("Declare queue error: %v. Retrying in %v...\n", err, retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("declare queue: %s error: %w", queueName, err)
	}

	for _, topic := range topics {
		for i := 0; i < maxRetries; i++ {
			err = r.ch.QueueBind(
				queue.Name, // queue name
				topic,      // routing key
				exchange,   // exchange
				false,
				nil,
			)
			if err == nil {
				break
			}
			r.logger.Printf("QueueBind error: %v. Retrying in %v...\n", err, retryInterval)
			time.Sleep(retryInterval)
		}

		if err != nil {
			return nil, fmt.Errorf("QueueBind error: %w", err)
		}
	}

	msgs, err := r.ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, fmt.Errorf("Consume error: %w", err)
	}

	out := make(chan amqp.Delivery)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		defer close(out)
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				out <- msg
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func (r *RabbitMQ) Close() {
	close(r.closeCh)
	r.wg.Wait()
	if err := r.ch.Close(); err != nil {
		r.logger.Printf("Failed to close channel: %s", err)
	}
	if err := r.conn.Close(); err != nil {
		r.logger.Printf("Failed to close connection: %s", err)
	}
}
