package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg []byte
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			delete(clients, conn)
			break
		}

		if len(msg) == 0 {
			fmt.Println("Received empty message, not broadcasting.")
			continue
		}

		fmt.Printf("Received: %s\n", msg)
		broadcastMessage(msg)
	}
}

func broadcastMessage(message []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error while writing message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/chat/", http.StripPrefix("/chat/", fs))

	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Запуск сервера на localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
