package main

import (
	"context"
	"fmt"
	"net/http"
	"webMessenger/socket"
	"webMessenger/user"
	"webMessenger/database"

	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/chat/", http.StatusPermanentRedirect)
}

func main() {
	database.ConnectDB()
	database.CreateIndexes()
	defer func() {
		if err := database.MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./resource"))
	router.PathPrefix("/chat/").Handler(http.StripPrefix("/chat/", fs))
	
	router.HandleFunc("/registration", user.Registration)
	router.HandleFunc("/ws", socket.GlobalChat)
	router.HandleFunc("/history", socket.GlobalHistory)
	router.HandleFunc("/chat", redirect)
	router.HandleFunc("/", redirect)
	router.HandleFunc("/chat/{user_name}", socket.PersonalChat)

	fmt.Println("Запуск сервера на localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
