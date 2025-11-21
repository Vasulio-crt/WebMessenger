package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"sync"
	
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message определяет структуру сообщения чата
type Message struct {
	SenderID string `json:"senderId"`
	Text     string `json:"text"`
}

const maxMessageLength = 500
const maxHistorySize = 50

const ResetColor string = "\033[0m"
const RedColor string = "\033[31m"
const GreenColor string = "\033[32m"
const YellowColor string = "\033[33m"

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = sync.Mutex{}

var messageHistory []Message
var historyMutex = sync.Mutex{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(YellowColor + "msg:  ", msg) // -----debug-----

		if len(msg.Text) > maxMessageLength {
			fmt.Printf("Message from %s rejected, too long: %d chars\n", msg.SenderID, len(msg.Text))
			continue
		}

		historyMutex.Lock()
		messageHistory = append(messageHistory, msg)

		fmt.Println("messageHistory:  ", messageHistory) // -----debug-----

		if len(messageHistory) > maxHistorySize {
			messageHistory = messageHistory[len(messageHistory) - maxHistorySize:]
		}
		historyMutex.Unlock()

		fmt.Printf("msgBytes: %s%s\n", bodyBytes, ResetColor) // -----debug-----
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

func handleHistory(w http.ResponseWriter, r *http.Request) {
	historyMutex.Lock()
	defer historyMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageHistory)
}

func main() {
	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/chat/", http.StripPrefix("/chat/", fs))

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/history", handleHistory)

	fmt.Println("Запуск сервера на localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
