package chats

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const maxMessageLength int = 1024

var upgrader = websocket.Upgrader{
	ReadBufferSize: maxMessageLength + 256,
	WriteBufferSize: maxMessageLength + 256,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	From string `json:"from"`
	Text string `json:"text"`
	Timestamp int32 `json:"timestamp"`
}

type MessageFromTo struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
	Timestamp int32 `json:"timestamp"`
	Type string `json:"type"`
}

type Info struct {
	Users    [2]string `bson:"users"`
	ChatName string    `bson:"chatName"`
}
