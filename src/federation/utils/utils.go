package utils

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "failed to establish connection to msg broker")
	ch, err := conn.Channel()
	failOnError(err, "failed to get channel for connection")
	return conn, ch
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// func to get queue based on routing key name and returns reference to queueu
func GetQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare queue")
	return &q
}
