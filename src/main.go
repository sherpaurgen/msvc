package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	server()
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	//direct xchange queue
	q, err := ch.QueueDeclare(
		"Q_1", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "failed to declare queue")
	return conn, ch, &q
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("this is msg samlpe..."),
	}
	//doesnt have specific name so its default exchange n
	ch.Publish("", q.Name, false, false, msg)
}
