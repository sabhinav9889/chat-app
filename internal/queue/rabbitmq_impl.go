package queue

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

// RabbitMQClient is an implementation of the QueueClient interface for RabbitMQ
type RabbitMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQClient creates a new instance of RabbitMQClient
func NewRabbitMQClient(amqpURL string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &RabbitMQClient{
		connection: conn,
		channel:    ch,
	}, nil
}

// PublishMessage sends a message to a RabbitMQ queue
func (r *RabbitMQClient) PublishMessage(ctx context.Context, queue string, message []byte) error {
	err := r.channel.Publish(
		"",    // exchange
		queue, // routing key (queue name)
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v", queue, err)
	}
	return nil
}

// ConsumeMessages listens for messages on a RabbitMQ queue and processes them
func (r *RabbitMQClient) ConsumeMessages(ctx context.Context, queue string, handler func(message []byte) error) error {
	msgs, err := r.channel.Consume(
		queue, // queue name
		"",    // consumer tag
		true,  // auto-acknowledge
		false, // exclusive
		false, // no-local
		true,  // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %v", err)
	}

	// Consume messages in a loop
	for msg := range msgs {
		err := handler(msg.Body)
		if err != nil {
			fmt.Printf("error handling message: %v\n", err)
		}
	}
	return nil
}

// Close gracefully shuts down the RabbitMQ client and channel
func (r *RabbitMQClient) Close() {
	r.channel.Close()
	r.connection.Close()
}
