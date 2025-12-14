package chats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"webMessenger/database"
	"webMessenger/user"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageGlobal struct {
	From string `json:"from"`
	Text string `json:"text"`
}

var clients = make(map[*websocket.Conn]string)
var clientsMutex sync.Mutex

func GlobalChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/registration", http.StatusFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Error GlobalChat: ", err)
	}
	defer conn.Close()

	userName, err := user.Get_user_name(r, cookie)
	if err != nil {
		http.Error(w, "User name not found", http.StatusInternalServerError)
		return
	}

	clientsMutex.Lock()
	clients[conn] = userName
	clientsMutex.Unlock()

	for {
		var msg MessageGlobal
		if err := conn.ReadJSON(&msg); err != nil {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			log.Println("Error ReadJSON:", err)
			break
		}

		if len(msg.Text) == 0 {
			continue
		}

		if len(msg.Text) > maxMessageLength {
			msg.From = "server"
			msg.Text = "Your message is too long (max 500 chars)"
			conn.WriteJSON(msg)
			continue
		}

		// Сохранение сообщения в БД
		collection := database.GetCollection("globalMessages")
		_, err = collection.InsertOne(context.Background(), msg)
		if err != nil {
			log.Println("Error inserting message to MongoDB:", err)
			continue
		}

		msg.From = userName
		broadcastMessage(&msg)
	}
}

func broadcastMessage(message *MessageGlobal) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Error while writing message:", err)
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
