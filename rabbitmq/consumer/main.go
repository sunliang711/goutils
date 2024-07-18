package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sunliang711/goutils/rabbitmq"
)

func main() {

	url := "amqp://guest:guest@10.1.9.120:5673/"

	rabbitMQ, err := rabbitmq.NewRabbitMQ(url, 5)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	rabbitMQ.AddConsumer("exchange001", "topic", []string{"topic1", "topic2"}, func(msg amqp091.Delivery) {
		fmt.Printf("got message: %v\n", string(msg.Body))
		msg.Ack(false)
	})
	err = rabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Waiting for exit signal")
		<-sigs
		rabbitMQ.Close()
		log.Println("Gracefully shutting down")
		os.Exit(0)
	}()

	// time.Sleep(time.Second * 10)
	// rabbitMQ.Close()
	select {}
}
