package main

import (
	"context"
	"fmt"
	"net/http"
	"webMessenger/chats"
	"webMessenger/database"
	"webMessenger/user"

	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/globalChat", http.StatusPermanentRedirect)
}

/* func HttpErrorPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/ErrorPage.html")
} */

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
	FS := http.FileServer(http.Dir("resource/"))

	router := mux.NewRouter()
	router.PathPrefix("/resource/").Handler(http.StripPrefix("/resource/", FS))

	// PersonalChat
	router.HandleFunc("/chat/{userName}", chats.GetChat)
	router.HandleFunc("/chat/find/{userName}", chats.FindChat)
	router.HandleFunc("/chat/{userName}/history", chats.GetChat)
	// router.HandleFunc("/wsp", chats.PersonalChat)

	router.HandleFunc("/registration", user.GetRegistration).Methods(http.MethodGet)
	router.HandleFunc("/login", user.GetLogin).Methods(http.MethodGet)
	router.HandleFunc("/register", user.Registration).Methods(http.MethodPost)
	router.HandleFunc("/login", user.Login).Methods(http.MethodPost)

	// GlobalChat
	router.HandleFunc("/globalChat", chats.GetGlobalChat)
	router.HandleFunc("/ws", chats.GlobalChat)
	router.HandleFunc("/history", chats.GlobalHistory)
	router.HandleFunc("/chat", redirect)
	router.HandleFunc("/", redirect)

	fmt.Println("Запуск сервера на localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
