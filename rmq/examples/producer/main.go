package main

import (
	"fmt"
	"log"
	"time"

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

	// 3. 添加生产者
	rabbitMQ.AddProducer("exchange001", exchangeOptions)

	// 4. 连接
	err = rabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	i := 0
	for {
		// 5. 发送消息
		log.Printf("Push message: %v\n", i)
		err = rabbitMQ.Publish("exchange001", "topic1", []byte(fmt.Sprintf("Hello, World: %v", i)))
		if err != nil {
			log.Printf("Failed to publish message: %s\n", err)
			time.Sleep(time.Second * 1)
			continue
		}

		// err = rabbitMQ.Publish("exchange001", "topic2", []byte(fmt.Sprintf("Hello, World: %v", i*10)))
		// if err != nil {
		// 	log.Fatalf("Failed to publish message: %s\n", err)
		// }
		i += 1
		time.Sleep(time.Second * 1)

	}
}
