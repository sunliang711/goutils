package rmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL          string
	ReconnectSec int // 重连间隔

	Producers map[string]ProducerConfig // 用于生产消息, key为exchange name, value为配置信息，包括: 类型
	Consumers map[string]ConsumerConfig // 用于消费信息, key为exchange name, value为配置信息，包括: 类型,topics,handler
}

type ProducerConfig struct {
	ExchangeOptions ExchangeOptions
	// TODO: other options
}

type ConsumerConfig struct {
	ExchangeOptions ExchangeOptions

	Handler        MessageHandlerFunc // 消息处理handler
	Topic          string             // type为topic|direct时的topics
	QueueOptions   QueueOptions
	ConsumeOptions ConsumeOptions
}

type ExchangeOptions struct {
	Type      string     // "topic|direct" "fanout"
	Durable   bool       // durable 如果设置为 true，交换机将在服务器重启后仍然存在
	AutoDel   bool       // auto-deleted 如果设置为 true，当没有队列绑定到交换机时，交换机会被自动删除。
	Internal  bool       // internal 如果设置为 true，交换机是内部使用的，不能被客户端直接发送消息。
	NoWait    bool       // no-wait 如果设置为 true，不等待服务器的确认。
	Arguments amqp.Table // arguments Table 类型，表示一个键值对的字典，用于指定交换机的额外参数。
}

type QueueOptions struct {
	Name      string     // name 队列名称
	Durable   bool       // durable 如果设置为 true，队列将在服务器重启后仍然存在。
	AutoDel   bool       // auto-deleted 如果设置为 true，当没有消费者使用时，队列会被自动删除。
	Exclusive bool       // exclusive 如果设置为 true，只有创建者可以使用的私有队列，断开后自动删除。
	NoWait    bool       // no-wait 如果设置为 true，不等待服务器的确认。
	Arguments amqp.Table // arguments Table 类型，表示一个键值对的字典，用于指定队列的额外参数。

	BindNoWait bool       // no-wait 如果设置为 true，不等待服务器的确认。
	BindArgs   amqp.Table // arguments Table 类型，表示一个键值对的字典，用于指定绑定的额外参数。
}

type ConsumeOptions struct {
	AutoAck   bool       // auto-ack 如果设置为 true，消息一旦被接收，客户端自动发送回执。
	Exclusive bool       // exclusive 如果设置为 true，队列会变成私有的，只有这个连接可以消费这个队列。
	NoLocal   bool       // no-local 如果设置为 true，表示不能将同一个connection中发送的消息传送给这个connection中的消费者。
	NoWait    bool       // no-wait 如果设置为 true，不等待服务器的确认。
	Arguments amqp.Table // arguments Table 类型，表示一个键值对的字典，用于指定消费者的额外参数。
}

type MessageHandlerFunc func(msg amqp.Delivery)

func messageHandler(ctx context.Context, msgs <-chan amqp.Delivery, handler MessageHandlerFunc) {
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
