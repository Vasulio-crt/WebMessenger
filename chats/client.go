package chats

import (
	"context"
	"log"
	"sync"
	"time"
	"webMessenger/database"
	"webMessenger/user/utilities"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	Name string
	Conn *websocket.Conn
}

func (c *Client) SendMessage(msg MessageFromTo) error {
	return c.Conn.WriteJSON(msg)
}

// ------- Hub -------
type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) AddClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.Name] = c
	log.Printf("Пользователь подключен: %s", c.Name)
}

func (h *Hub) RemoveClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.clients[id]; ok {
		c.Conn.Close()
		delete(h.clients, id)
		log.Printf("Пользователь отключен: %s", id)
	}
}

// SendPrivateMessage ищет получателя и пишет ему
func (h *Hub) SendPrivateMessage(msg MessageFromTo) {
	collectionInfo := database.GetCollectionHistory("info")
	filter := bson.M{
		"users": bson.M{
			"$all": bson.A{msg.From, msg.To},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var info Info
	var hashChat string
	if err := collectionInfo.FindOne(ctx, filter).Decode(&info); err != nil {
		if err == mongo.ErrNoDocuments {
			hashChat = utilities.HashString(msg.From + msg.To)
			_, err = collectionInfo.InsertOne(ctx, bson.M{"users": bson.A{msg.From, msg.To}, "chatName": hashChat})
			if err != nil {
				log.Println("Error inserting message to MongoDB:", err)
				return
			}
		}
		log.Fatal(err)
	} else {
		hashChat = info.ChatName
	}

	h.mu.RLock()
	recipient, ok := h.clients[msg.To]
	h.mu.RUnlock()

	if ok {
		err := recipient.SendMessage(msg)
		if err != nil {
			log.Printf("Ошибка отправки пользователю %s: %v", msg.To, err)
			h.RemoveClient(recipient.Name)
		}
	} else {
		collection := database.GetCollection(hashChat)
		msgHistory := Message{
			From: msg.From,
			Text: msg.Text,
		}
		if _, err := collection.InsertOne(ctx, msgHistory); err != nil {
			log.Println("Error inserting message to MongoDB:", err)
			return
		}
	}
}
