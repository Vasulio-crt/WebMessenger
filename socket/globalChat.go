package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"webMessenger/database"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	SenderID  string    `json:"senderId"`
	Text      string    `json:"text"`
}

const maxMessageLength int = 500 * 2

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = sync.Mutex{}

func GlobalChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	for {
		_, bodyBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}

		if len(bodyBytes) == 0 {
			fmt.Println("Received empty message, not broadcasting.")
			continue
		}

		var msg Message
		if err := json.Unmarshal(bodyBytes, &msg); err != nil {
			fmt.Println("Error unmarshaling message:", err)
			continue
		}

		if len(msg.Text) > maxMessageLength {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"senderId":"server","text":"Your message is too long (max 500 chars)"}`))
			continue
		}

		// Сохранение сообщения в БД
		collection := database.GetCollection("globalMessages")
		_, err = collection.InsertOne(context.Background(), msg)
		if err != nil {
			fmt.Println("Error inserting message to MongoDB:", err)
			continue
		}

		broadcastMessage(bodyBytes)
	}
}

func broadcastMessage(message []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error while writing message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func GlobalHistory(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("globalMessages")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})
	findOptions.SetLimit(50)

	cursor, err := collection.Find(r.Context(), bson.D{}, findOptions)
	if err != nil {
		http.Error(w, "Failed to retrieve message history", http.StatusInternalServerError)
		fmt.Println("Error finding messages in MongoDB:", err)
		return
	}
	defer cursor.Close(r.Context())

	var messageHistory []Message
	if err = cursor.All(r.Context(), &messageHistory); err != nil {
		http.Error(w, "Failed to decode message history", http.StatusInternalServerError)
		fmt.Println("Error decoding messages from MongoDB:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageHistory)
}
