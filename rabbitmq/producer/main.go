package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sunliang711/goutils/rabbitmq"
)

func main() {

	url := "amqp://guest:guest@10.1.9.120:5673/"
	rabbitMQ, err := rabbitmq.NewRabbitMQ(url, 5)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	rabbitMQ.AddProducer("exchange001", "topic", []string{"topic1", "topic2"})
	err = rabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	i := 0
	for {
		// Publish example
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
