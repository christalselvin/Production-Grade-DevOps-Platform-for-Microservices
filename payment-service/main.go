package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://admin:admin123@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	q, err := ch.QueueDeclare(
		"orderQueue",
		true,
		false,
		false,
		false,
		nil,
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	fmt.Println("Payment service waiting for orders...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			fmt.Println("Payment processing order:", string(d.Body))

		}
	}()

	<-forever
}
