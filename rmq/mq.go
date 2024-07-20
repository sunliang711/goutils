package rmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
			URL:          url,
			ReconnectSec: reconnect,
			Consumers:    make(map[string]ConsumerConfig),
			Producers:    make(map[string]ProducerConfig),
		},
		// closeCh:    make(chan struct{}),
		logger:     log.New(os.Stdout, "|RMQ| ", log.LstdFlags),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	return r, nil
}

func (r *RabbitMQ) AddConsumer(exchangeName string, topic string, handler MessageHandlerFunc, exchangeOptions ExchangeOptions, queueOptions QueueOptions, consumeOptions ConsumeOptions) error {
	config := ConsumerConfig{
		ExchangeOptions: exchangeOptions,
		Handler:         handler,
		Topic:           topic,
		QueueOptions:    queueOptions,
		ConsumeOptions:  consumeOptions,
	}

	if handler == nil {
		return fmt.Errorf("exchange %s handler is nil", exchangeName)
	}

	r.config.Consumers[exchangeName] = config
	return nil
}

func (r *RabbitMQ) AddProducer(exchangeName string, exchangeOptions ExchangeOptions) error {
	exchange := ProducerConfig{
		ExchangeOptions: exchangeOptions,
	}

	r.config.Producers[exchangeName] = exchange
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
	for exchangeName, producer := range r.config.Producers {
		r.logger.Printf("Declare exchange: %s\n", exchangeName)
		err = r.ch.ExchangeDeclare(
			exchangeName,                       // name
			producer.ExchangeOptions.Type,      // type
			producer.ExchangeOptions.Durable,   // durable
			producer.ExchangeOptions.AutoDel,   // auto-deleted
			producer.ExchangeOptions.Internal,  // internal
			producer.ExchangeOptions.NoWait,    // no-wait
			producer.ExchangeOptions.Arguments, // arguments
		)
		if err != nil {
			return err
		}

	}

	// 声明consume exchange
	for exchangeName, consumer := range r.config.Consumers {
		r.logger.Printf("Declare exchange: %s\n", exchangeName)
		err = r.ch.ExchangeDeclare(
			exchangeName,                       // name
			consumer.ExchangeOptions.Type,      // type
			consumer.ExchangeOptions.Durable,   // durable
			consumer.ExchangeOptions.AutoDel,   // auto-deleted
			consumer.ExchangeOptions.Internal,  // internal
			consumer.ExchangeOptions.NoWait,    // no-wait
			consumer.ExchangeOptions.Arguments, // arguments
		)
		if err != nil {
			return err
		}
	}

	r.handleReconnect()

	// 遍历consume exchange, 声明queue并绑定exchange
	for exchangeName, consumer := range r.config.Consumers {
		ch, err := r.consume(r.ctx, exchangeName, &consumer)
		if err != nil {
			return err
		}
		go messageHandler(r.ctx, ch, consumer.Handler)
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

func (r *RabbitMQ) consume(ctx context.Context, exchangeName string, consumerConfig *ConsumerConfig) (<-chan amqp.Delivery, error) {
	r.consumeMux.Lock()
	defer r.consumeMux.Unlock()

	var queue amqp.Queue
	var err error
	maxRetries := 5
	retryInterval := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		r.logger.Printf("Declare queue: %s\n", consumerConfig.QueueOptions.Name)
		queue, err = r.ch.QueueDeclare(
			consumerConfig.QueueOptions.Name,      // name
			consumerConfig.QueueOptions.Durable,   // durable
			consumerConfig.QueueOptions.AutoDel,   // delete when unused
			consumerConfig.QueueOptions.Exclusive, // exclusive
			consumerConfig.QueueOptions.NoWait,    // no-wait
			consumerConfig.QueueOptions.Arguments, // arguments
		)
		if err == nil {
			break
		}
		r.logger.Printf("Declare queue error: %v. Retrying in %v...\n", err, retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("declare queue: %s error: %w", consumerConfig.QueueOptions.Name, err)
	}

	// 绑定exchange和queue
	for i := 0; i < maxRetries; i++ {
		err = r.ch.QueueBind(
			queue.Name,           // queue name
			consumerConfig.Topic, // routing key
			exchangeName,         // exchange
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

	msgs, err := r.ch.Consume(
		queue.Name,                              // queue
		"",                                      // consumer
		consumerConfig.ConsumeOptions.AutoAck,   // auto-ack
		consumerConfig.ConsumeOptions.Exclusive, // exclusive
		consumerConfig.ConsumeOptions.NoLocal,   // no-local
		consumerConfig.ConsumeOptions.NoWait,    // no-wait
		consumerConfig.ConsumeOptions.Arguments, // args
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
			r.logger.Printf("Consume %s stopped\n", queue.Name)
		}()
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				out <- msg
			case <-ctx.Done():
				r.logger.Printf("Consume %s done\n", queue.Name)
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
