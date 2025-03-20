package queue

import "context"

// QueueClient defines methods for interacting with a message queue (RabbitMQ)
type QueueClient interface {
    PublishMessage(ctx context.Context, queue string, message []byte) error
    ConsumeMessages(ctx context.Context, queue string, handler func(message []byte) error) error
}
