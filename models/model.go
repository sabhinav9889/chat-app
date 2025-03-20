package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	UserId   string
	IsSeeker bool
	IsClose  bool
	Conn     *websocket.Conn
}

type Message struct {
	MessageId    string   `json:"message_id"`
	Type         string   `json:"type"`               // Message type (e.g., "status", "content", "file")
	Content      string   `json:"content"`            // The message content (e.g., status message or text)
	Status       string   `json:"status,omitempty"`   // Optional: status of the user (e.g., "online", "typing")
	FileName     string   `json:"filename,omitempty"` // Optional: name of the file
	FileType     string   `json:"filetype,omitempty"` // Optional: MIME type of the file (e.g., "application/pdf")
	FileData     []byte   `json:"filedata,omitempty"` // Optional: file data (binary content)
	ReceiverList []string `json:"receiver_id"`        // Receiver who receive this message from sender
	SenderID     string   `json:"sender_id"`          // sender id can send messages to the multiple users
	TimeStamp    int64    `json:"timestamp"`          // timestamp of message received
}

type SendMessage struct {
	MessageId     string `json:"message_id"`
	MessageStatus string `json:"message_status"` // edit or not
	Content       string `json:"content"`        // The message content (e.g., status message or text)
	ReceiverID    string `json:"receiver_id"`    // Receiver who receive this message from sender
	SenderID      string `json:"sender_id"`      // sender id can send messages to the multiple users
	TimeStamp     int64  `json:"timestamp"`      // timestamp of message received
}

type LastSeenMessage struct {
	ReceiverID string `json:"receiver_id"` // Receiver's id
	IsOnline   bool   `json:"is_online"`
	LastSeen   int    `json:"last_seen"` // last time when receiverid was online
}

type ChatResponse struct {
	MessageID string    `json:"MessageID"`
	Sender    string    `json:"Sender"`
	Receiver  string    `json:"Receiver"`
	Message   string    `json:"Message"`
	MsgType   string    `json:"MsgType"`
	Timestamp time.Time `json:"Timestamp"`
	Status    string    `json:"Status"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type ChatMessagesResponse struct {
	Sender   string         `json:"sender"`
	Receiver string         `json:"receiver"`
	Messages []ChatResponse `json:"messages"`
}

type ChatHistoryResponse struct {
	Sender   string         `json:"sender"`
	Receiver string         `json:"receiver"`
	Messages []ChatResponse `json:"messages"`
}
