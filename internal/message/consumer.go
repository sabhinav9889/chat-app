package message

import (
	"chat-app/internal/queue"
	internal_websocket "chat-app/internal/websocket"
	"context"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

// ConsumeMessagesFromQueue listens for messages from RabbitMQ and forwards them to WebSocket clients
func ConsumeMessagesFromQueue(conn *websocket.Conn, wsClient *internal_websocket.WebSocketClientImpl) {
	// Initialize the RabbitMQ client
	rabbitMQClient, err := queue.NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer rabbitMQClient.Close()

	// Start consuming messages from the "chatQueue"
	err = rabbitMQClient.ConsumeMessages(context.Background(), "chatQueue", func(message []byte) error {
		// Print the consumed message
		fmt.Printf("Consumed message: %s\n", message)

		// Send the message to the WebSocket client
		if err := wsClient.SendMessage(conn, message); err != nil {
			log.Printf("Error sending message to WebSocket: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Error consuming messages: %v", err)
	}
}
