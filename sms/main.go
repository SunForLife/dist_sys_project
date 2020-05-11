package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	time.Sleep(45 * time.Second)

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("sms", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	for {
		time.Sleep(time.Second)
		msgs, err := ch.Consume(
			q.Name, "", true, false, false, false, nil)
		failOnError(err, "Failed to register a consumer")
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}
}
