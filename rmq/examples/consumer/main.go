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

	caCertBytes := []byte(`-----BEGIN CERTIFICATE-----
MIIDtjCCAp6gAwIBAgITCeBKIs7Ew4XSsy/Wul3kCQ9KsTANBgkqhkiG9w0BAQsF
ADBrMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
U2FuIEZyYW5jaXNjbzEOMAwGA1UECgwFTXlPcmcxDzANBgNVBAsMBk15RGVwdDEO
MAwGA1UEAwwFbXktY2EwHhcNMjQwNzIyMTA0MTAzWhcNMjUwNzIyMTA0MTAzWjBr
MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2Fu
IEZyYW5jaXNjbzEOMAwGA1UECgwFTXlPcmcxDzANBgNVBAsMBk15RGVwdDEOMAwG
A1UEAwwFbXktY2EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDXiqz/
sS/KeoL4nVmgZjYNO9mZpZz2dVA5C5Ep8E7HFF93/ctrsdukZ+MMi5u4rfhDO6xn
dL9FV7A/IwYWBNXheuN5DDhSIyAu8ZTRHzEFYpQ4l2zLqIb04FaONCO1LReh30SX
qv97n7BihiqqVnpQLpriScKiVk4Cw21LyH6BxkzXqgMl0HmvJ8mg1IQlqkTQhPHn
ctsB91uqyiShPU841bqskkpwSXXdLOVbBIQ8IMumTCOoWdlckvqcnfZsrVcLGYGa
vNuaUbEFreMNpYN8r8nVMJCQIpUdqN8I7KFDAxOqgI/StwNFrWVEdyAlMUpMJKZ4
IGPAaRkAXxr7aEJZAgMBAAGjUzBRMB0GA1UdDgQWBBSq/ym/7JkJW/+BGwHFsEch
9+oTsjAfBgNVHSMEGDAWgBSq/ym/7JkJW/+BGwHFsEch9+oTsjAPBgNVHRMBAf8E
BTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQDDXdFxJJBKuCIYoCyY/NTR7lok3Hk6
CvT97LMztoiIciR7MJqNK7e1xgFFs5DfSgYqdzSdUjvUG+WoUk/jXJzSsanKsErF
cfLXMlKig2EAlBgACfgyRuHghTMq2YoF4wSkpJJUBnjejmAknYtGy2fU2H7rsTOK
BNT0wMUQ10C1eJnYoR7QgyWw3zZOAdo1ivBKagMvdyl+Bz99FrvGQVj4F0EYrvam
7Ee+lNesGfJm6CHuN9eC2uxWNJCxCgr0C+lhnIW6c/23CtbUTqC79ZbBpSqxLUth
AsBLPhYADl31nkFiJKO/r7Ewdwhwahl3lkCxnigBAeUnm0VICtmlNBdu
-----END CERTIFICATE-----`)

	url := "amqps://guest:guest@10.1.9.66:5671/"
	// 1. 构建实例
	rabbitMQ, err := rmq.NewRabbitMQ(url, 5, caCertBytes, nil, nil)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// 2. 队列参数
	queueOptions := rmq.QueueOptions{
		Name:    "qq",
		Durable: true,
	}
	// 3. 消费参数
	consumeOptions := rmq.ConsumeOptions{}
	// 4. 添加生产者
	exchangeName := "exchange001"
	rabbitMQ.AddConsumer(exchangeName, "topic1", func(msg amqp091.Delivery) {
		fmt.Printf("got message: %v\n", string(msg.Body))
		msg.Ack(false)
	}, queueOptions, consumeOptions)

	// 5. 连接
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
