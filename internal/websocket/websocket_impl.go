package websocket

import (
	"chat-app/constants"
	"chat-app/internal/queue"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// WebSocketClientImpl is an implementation of the WebSocketClient interface
type WebSocketClientImpl struct{}

// NewWebSocketClient creates a new instance of WebSocketClientImpl
func NewWebSocketClient() *WebSocketClientImpl {
	return &WebSocketClientImpl{}
}

// UpgradeConnection upgrades an HTTP connection to a WebSocket connection
func (ws *WebSocketClientImpl) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade to websocket: %v", err)
	}
	return conn, nil
}

// SendMessage sends a message to the WebSocket connection
func (ws *WebSocketClientImpl) SendMessage(conn *websocket.Conn, message []byte) error {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// Read incomming message to the websocket connection
func (ws *WebSocketClientImpl) ReadMessage(conn *websocket.Conn, UserId string, qu *queue.RabbitMQClient) {
	ctx := context.Background()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Err(err).Msg("Failed to read message")
			break
		}
		err = qu.PublishMessage(ctx, constants.QueueName, msg)
		if err != nil {
			log.Err(err).Msg("Fail to publish message in queue")
			continue
		}
		log.Info().Msg("Message published successfully in queue")
	}
}

// CloseConnection closes the WebSocket connection
func (ws *WebSocketClientImpl) CloseConnection(conn *websocket.Conn) error {
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}
