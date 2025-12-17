package chats

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const maxMessageLength int = 500

type Message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type MessageHistory struct {
	From string `json:"from"`
	Text string `json:"text"`
}

type Info struct {
	Users    [2]string `bson:"users"`
	ChatName string    `bson:"chatName"`
}
