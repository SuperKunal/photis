package services

import (
	"github.com/streadway/amqp"
	"log"
)

type Client struct {
	conn *amqp.Connection
	ch *amqp.Channel
	q *amqp.Queue
}

func NewRabbitMqClient(addr, queueName string) *Client {
	conn, err := amqp.Dial(addr)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &Client{
		conn: conn,
		ch: ch,
		q: &q,
	}
}

func (c *Client) Publish(content string) {
	err := c.ch.Publish(
		"",     // exchange
		c.q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}