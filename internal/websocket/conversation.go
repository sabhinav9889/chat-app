package websocket

import (
	"chat-app/constants"
	"chat-app/internal/cache"
	"chat-app/internal/queue"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rs/zerolog/log"
)

var webskt *WebSocketClientImpl
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; adjust this as per your securit	y needs
	},
}

func StartChat() {
	webskt = NewWebSocketClient()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := webskt.UpgradeConnection(w, r)
		if err != nil {
			log.Err(err).Msg("Error building websocket")
			return
		}
		conn.SetReadDeadline(time.Now().Add(3600 * time.Second))
		userId := r.URL.Query().Get("user_id")
		queue, err := queue.NewRabbitMQClient(constants.AmqpURL)
		redis := cache.NewRedisCacheClient(constants.RedisAdd)
		redis.Set(context.Background(), userId, conn)
		webskt.ReadMessage(conn, userId, queue)
	})
	log.Info().Msg("Chat Server started on : 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Err(err).Msg("Error starting server: ")
	}
}
