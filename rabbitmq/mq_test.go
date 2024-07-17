package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestMQ(t *testing.T) {
	config := Config{
		URL:          "amqp://guest:guest@10.1.9.120:5672/",
		ReconnectSec: 5,
		Exchanges: map[string]ExchangeConfig{
			"exchange001": {
				Type:   "topic",
				Topics: []string{"topic1", "topic2"},
			},
			"exchange002": {
				Type:   "topic",
				Topics: []string{"topic3", "topic4"},
			},
		},
	}

	rabbitMQ, err := NewRabbitMQ(config)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
		rabbitMQ.Close()
		log.Println("Gracefully shutting down")
		os.Exit(0)
	}()

	go func() {
		// time.Sleep(time.Second * 2)
		log.Printf("Consume1\n")

		msgs, err := rabbitMQ.Consume(ctx, "exchange001", "queue1", []string{"topic1", "topic2"})
		if err != nil {
			log.Fatalf("Failed to consume messages: %s", err)
		}
		for msg := range msgs {
			log.Printf("Received message: %s\n", msg.Body)
		}

	}()

	go func() {
		// time.Sleep(time.Second * 2)
		log.Printf("Consume2\n")
		msgs, err := rabbitMQ.Consume(ctx, "exchange002", "queue2", []string{"topic3", "topic4"})
		if err != nil {
			log.Fatalf("Failed to consume messages: %s", err)
		}

		for msg := range msgs {
			log.Printf("Received message: %s\n", msg.Body)
		}
	}()

	i := 0
	for {
		// Publish example
		log.Printf("Push message: %v\n", i)
		err = rabbitMQ.Publish("exchange001", "topic1", []byte(fmt.Sprintf("Hello, World: %v", i)))
		if err != nil {
			log.Fatalf("Failed to publish message: %s\n", err)
		}

		err = rabbitMQ.Publish("exchange001", "topic2", []byte(fmt.Sprintf("Hello, World: %v", i*10)))
		if err != nil {
			log.Fatalf("Failed to publish message: %s\n", err)
		}
		i += 1
		time.Sleep(time.Second * 2)
	}
}
