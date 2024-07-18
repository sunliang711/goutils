package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL string
	// Exchanges        map[string]ExchangeConfig
	ReconnectSec int // 重连间隔

	ProduceExchanges map[string]ExchangeConfig // 用于生产消息, key为exchange name, value为配置信息，包括: 类型
	ConsumeExchanges map[string]ExchangeConfig // 用于消费信息, key为exchange name, value为配置信息，包括: 类型,topics,handler
}

// context.Context用于取消操作
// type MessageHandler func(context.Context, <-chan amqp.Delivery)

type MessageHandlerFunc func(msg amqp.Delivery)

func messageHandler(ctx context.Context, msgs <-chan amqp091.Delivery, handler MessageHandlerFunc) {
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("Channel closed\n")
				return
			}
			handler(msg)
		case <-ctx.Done():
			log.Printf("done,quit")
			return
		}
	}
}

type ExchangeConfig struct {
	Type    string             //  "topic|direct" "fanout" // TODO: topics must be empty when fanout
	Handler MessageHandlerFunc // 消息处理handler

	// Durable   bool       // durable
	// AutoDel   bool       // auto-deleted
	// Internal  bool       // internal
	// NoWait    bool       // no-wait
	// Arguments amqp.Table // arguments

	Topics []string // type为topic|direct时的topics
}

type RabbitMQ struct {
	config       Config
	conn         *amqp.Connection
	ch           *amqp.Channel
	wg           sync.WaitGroup
	reconnectMux sync.Mutex
	consumeMux   sync.Mutex

	logger *log.Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewRabbitMQ(url string, reconnect int) (*RabbitMQ, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &RabbitMQ{
		config: Config{
			URL:              url,
			ReconnectSec:     reconnect,
			ConsumeExchanges: make(map[string]ExchangeConfig),
			ProduceExchanges: make(map[string]ExchangeConfig),
		},
		// closeCh:    make(chan struct{}),
		logger:     log.New(os.Stdout, "|RMQ| ", log.LstdFlags),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	return r, nil
}

func (r *RabbitMQ) AddConsumer(exchangeName string, typ string, topics []string, handler MessageHandlerFunc) error {
	exchange := ExchangeConfig{
		Type:    typ,
		Handler: handler,
		Topics:  topics,
	}
	if exchange.Handler == nil {
		return fmt.Errorf("exchange %s handler is nil", exchangeName)
	}

	r.config.ConsumeExchanges[exchangeName] = exchange
	return nil
}

func (r *RabbitMQ) AddProducer(exchangeName string, typ string, topics []string) error {
	exchange := ExchangeConfig{
		Type:   typ,
		Topics: topics,
	}

	r.config.ProduceExchanges[exchangeName] = exchange
	return nil
}

func (r *RabbitMQ) Connect() error {
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

	// 声明produce exchange
	for name, config := range r.config.ProduceExchanges {
		r.logger.Printf("Declare exchange: %s\n", name)
		err = r.ch.ExchangeDeclare(
			name,        // name
			config.Type, // type
			true,        // durable
			false,       // auto-deleted
			false,       // internal
			false,       // no-wait
			nil,         // arguments
		)
		if err != nil {
			return err
		}

	}

	// 声明consume exchange
	for name, config := range r.config.ConsumeExchanges {
		r.logger.Printf("Declare exchange: %s\n", name)
		err = r.ch.ExchangeDeclare(
			name,        // name
			config.Type, // type
			true,        // durable
			false,       // auto-deleted
			false,       // internal
			false,       // no-wait
			nil,         // arguments
		)
		if err != nil {
			return err
		}
	}

	r.handleReconnect()

	// 遍历consume exchange, 声明queue并绑定exchange
	for name, exchange := range r.config.ConsumeExchanges {
		queue := fmt.Sprintf("queue_%s", name)
		ch, err := r.consume(r.ctx, name, queue, exchange.Topics)
		if err != nil {
			return err
		}
		go messageHandler(r.ctx, ch, exchange.Handler)
	}

	return nil
}

func (r *RabbitMQ) reconnect() {
	for {
		r.logger.Println("Attempting to reconnect...")
		err := r.Connect()
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
		err := <-r.conn.NotifyClose(make(chan *amqp.Error))
		if err != nil {
			r.logger.Printf("Connection closed: %s", err)
			r.reconnect()
		}
		// for {
		// 	err := <-r.conn.NotifyClose(make(chan *amqp.Error))
		// 	if err != nil {
		// 		r.logger.Printf("Connection closed: %s", err)
		// 		r.reconnect()
		// 	}
		// 	// Exit the goroutine if the connection is closed
		// 	if r.conn.IsClosed() {
		// 		return
		// 	}
		// }
	}()
}

// TODO: mandatory immediate
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

func (r *RabbitMQ) consume(ctx context.Context, exchange, queueName string, topics []string) (<-chan amqp.Delivery, error) {
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
		false,      // auto-ack
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
		defer func() {
			defer r.wg.Done()
			defer close(out)
			r.logger.Printf("Consume %s stopped\n", queueName)
		}()
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				out <- msg
			case <-ctx.Done():
				r.logger.Printf("Consume %s done\n", queueName)
				return
			}
		}
	}()
	return out, nil
}

func (r *RabbitMQ) Close() {
	r.logger.Println("Closing RabbitMQ connection...")
	r.cancelFunc()
	r.wg.Wait()
	r.logger.Println("All consumers stopped")
	if err := r.ch.Close(); err != nil {
		r.logger.Printf("Failed to close channel: %s", err)
	}
	if err := r.conn.Close(); err != nil {
		r.logger.Printf("Failed to close connection: %s", err)
	}
}
