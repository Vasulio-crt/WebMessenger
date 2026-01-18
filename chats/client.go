package chats

import (
	"context"
	"errors"
	"log"
	"sort"
	"sync"
	"webMessenger/database"
	"webMessenger/user/utilities"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

type Hub struct {
	clients map[string]*websocket.Conn //Name: Conn
	mu      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) AddClient(name string, conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[name] = conn
	h.mu.Unlock()
}

func (h *Hub) RemoveClient(name string) {
	h.mu.Lock()
	delete(h.clients, name)
	h.mu.Unlock()
}

// SendPrivateMessage ищет получателя и пишет ему
func (h *Hub) SendPrivateMessage(msg MessageFromTo) error {
	users := []string{msg.From, msg.To}
	sort.Strings(users)
	hashChat := utilities.HashString(users[0] + users[1])

	h.mu.RLock()
	recipient, ok := h.clients[msg.To]
	h.mu.RUnlock()

	if ok {
		err := recipient.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending to user %s: %v. Deleting client.\n", msg.To, err)
			h.RemoveClient(msg.To)
		}
	}
	// Сохранение истории чата
	collectionChats := database.GetCollectionHistory(hashChat)
	bsonM := bson.M{
		"From": msg.From,
		"Text": msg.Text,
	}
	if _, err := collectionChats.InsertOne(context.TODO(), bsonM); err != nil {
		return errors.New("Fail insertHistory")
	}
	return nil
}
