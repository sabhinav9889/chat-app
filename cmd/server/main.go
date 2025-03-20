package main

import (
	"fmt"
	"log"
	"net/http"
	"chat-app/internal/websocket"
	"chat-app/internal/message"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Create a new WebSocket client
	wsClient := websocket.NewWebSocketClient()

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsClient.UpgradeConnection(w, r)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer conn.Close()

	// Send a welcome message to the WebSocket client
	welcomeMessage := []byte("Welcome to the chat server!")
	if err := wsClient.SendMessage(conn, welcomeMessage); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	// Start consuming messages from RabbitMQ and forward them to the WebSocket client
	go message.ConsumeMessagesFromQueue(conn, wsClient)

	// Handle incoming messages from the WebSocket client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		// Handle the incoming message (publish to RabbitMQ and send echo)
		if err := message.HandleIncomingMessage(msg, conn, wsClient); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

func main() {
	// Handle WebSocket upgrade requests
	http.HandleFunc("/ws", handleWebSocket)

	// Start the WebSocket server on port 8080
	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
