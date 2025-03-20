package message

import (
	"chat-app/internal/queue"
	internal_websocket "chat-app/internal/websocket"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

// HandleIncomingMessage handles an incoming message and publishes it to RabbitMQ
func HandleIncomingMessage(message []byte, conn *websocket.Conn, wsClient *internal_websocket.WebSocketClientImpl) error {
	// Print the received message to the console (for debugging)
	fmt.Printf("Received message: %s\n", message)

	// Publish message to RabbitMQ queue
	err := publishMessageToQueue(message)
	if err != nil {
		log.Printf("Error publishing message to queue: %v", err)
		return err
	}

	// Optionally, send the message to the WebSocket client (echo)
	err = wsClient.SendMessage(conn, message)
	if err != nil {
		log.Printf("Error sending echo: %v", err)
		return err
	}

	return nil
}

// publishMessageToQueue publishes a message to the RabbitMQ queue
func publishMessageToQueue(message []byte) error {
	// Initialize RabbitMQ client
	rabbitMQClient, err := queue.NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ client: %v", err)
	}
	defer rabbitMQClient.Close()

	// Publish the message to the "chatQueue"
	err = rabbitMQClient.PublishMessage(context.Background(), "chatQueue", message)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	fmt.Println("Message published to RabbitMQ!")
	return nil
}
