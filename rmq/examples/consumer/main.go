package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sunliang711/goutils/rmq"
)

func main() {

	url := "amqp://guest:guest@10.1.9.66:5672/"

	// 1. 构建实例
	rabbitMQ, err := rmq.NewRabbitMQ(url, 5)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// 2. exchange参数
	exchangeOptions := rmq.ExchangeOptions{
		Type:    "topic",
		Durable: true,
	}
	// 3. 队列参数
	queueOptions := rmq.QueueOptions{
		Name:    "qq",
		Durable: true,
	}
	// 4. 消费参数
	consumeOptions := rmq.ConsumeOptions{}
	// 5. 添加生产者
	rabbitMQ.AddConsumer("exchange001", []string{"topic1"}, func(msg amqp091.Delivery) {
		fmt.Printf("got message: %v\n", string(msg.Body))
		msg.Ack(false)
	}, exchangeOptions, queueOptions, consumeOptions)

	// 6. 连接
	err = rabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Waiting for exit signal")
		<-sigs
		// 7. 关闭连接
		rabbitMQ.Close()
		log.Println("Gracefully shutting down")
		os.Exit(0)
	}()

	// time.Sleep(time.Second * 10)
	// rabbitMQ.Close()
	select {}
}
