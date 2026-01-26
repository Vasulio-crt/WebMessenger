package chats

import (
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

var clients = make(map[string]*websocket.Conn, 0)
var clientsMutex sync.Mutex

func GetGlobalChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || !user.Session_check(cookie){
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.ServeFile(w, r, "./pages/chat.html")
}

func GlobalChatWS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Error GlobalChat: ", err)
	}
	defer conn.Close()

	userName := user.Get_user_name(cookie)
	if userName == "" {
		http.Error(w, "User name not found", http.StatusInternalServerError)
		return
	}

	clientsMutex.Lock()
	clients[userName] = conn
	clientsMutex.Unlock()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			clientsMutex.Lock()
			conn.Close()
			delete(clients, userName)
			clientsMutex.Unlock()
			// log.Println("Error ReadJSON:", err)
			break
		}

		// Удаление сообщения с БД
		if msg.Text == "" {
			collection := database.GetCollection("globalMessages")
			_, err := collection.DeleteOne(r.Context(), bson.M{"timestamp": msg.Timestamp, "from": msg.From})
			if err != nil{
				log.Println("Error deleting message from MongoDB:", err)
				continue
			}
			continue
		}

		if len(msg.Text) > maxMessageLength {
			msg.From = "server"
			msg.Text = fmt.Sprintf("Your message is too long (max %d chars)", maxMessageLength)
			conn.WriteJSON(msg)
			continue
		}

		// Сохранение сообщения в БД
		collection := database.GetCollection("globalMessages")
		_, err = collection.InsertOne(r.Context(), msg)
		if err != nil {
			log.Println("Error inserting message to MongoDB:", err)
			continue
		}

		msg.From = userName
		broadcastMessage(&msg)
	}
}

func broadcastMessage(message *Message) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client, conn := range clients {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println("Error while writing message:", err)
			conn.Close()
			delete(clients, client)
		}
	}
}

func GlobalHistory(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("globalMessages")

	cursor, err := collection.Find(r.Context(), bson.D{}, options.Find().SetLimit(100))
	if err != nil {
		http.Error(w, "Failed to retrieve message history", http.StatusInternalServerError)
		fmt.Println("Error finding messages in MongoDB:", err)
		return
	}
	defer cursor.Close(r.Context())

	var chatHistory []Message
	if err = cursor.All(r.Context(), &chatHistory); err != nil {
		http.Error(w, "Failed to decode message history", http.StatusInternalServerError)
		fmt.Println("Error decoding messages from MongoDB:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatHistory)
}
