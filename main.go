package main

import (
	"context"
	"fmt"
	"net/http"
	"webMessenger/database"
	"webMessenger/chats"
	"webMessenger/user"

	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/globalChat/", http.StatusPermanentRedirect)
}

/* func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/registration", http.StatusFound)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// TODO: Проверить значение cookie в базе данных
		_ = cookie.Value

		next.ServeHTTP(w, r)
	})
} */

func main() {
	database.ConnectDB()
	defer func() {
		if err := database.MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()

	router.PathPrefix("/globalChat/").Handler(http.StripPrefix("/globalChat/", http.FileServer(http.Dir("./resource/globalChat"))))
	router.HandleFunc("/chat/{userName:[^.]+}", chats.PersonalChat)

	router.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./resource/registration/index.html")
	})
	router.HandleFunc("/register", user.Registration).Methods(http.MethodPost)
	router.HandleFunc("/ws", chats.GlobalChat)
	router.HandleFunc("/history", chats.GlobalHistory)
	router.HandleFunc("/chat", redirect)
	router.HandleFunc("/", redirect)

	fmt.Println("Запуск сервера на localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
