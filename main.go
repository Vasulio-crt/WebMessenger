package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webMessenger/chats"
	"webMessenger/database"
	"webMessenger/user"

	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/globalChat", http.StatusPermanentRedirect)
}

func authMiddleware(next http.Handler) http.Handler {
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

		isValid := user.Session_check(cookie)
		if !isValid {
			http.Redirect(w, r, "/registration", http.StatusFound)
			return
		}

		userName := user.Get_user_name(cookie)
		ctx := context.WithValue(r.Context(), "userName", userName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	database.ConnectDB()
	if err := user.Init_sessions(); err != nil {
		log.Fatalln("user.Init_sessions error:", err)
	}
	FS := http.FileServer(http.Dir("resource/"))

	router := mux.NewRouter()
	router.PathPrefix("/resource/").Handler(http.StripPrefix("/resource/", FS))

	// User
	router.HandleFunc("/registration", user.GetRegistration).Methods(http.MethodGet)
	router.HandleFunc("/login", user.GetLogin).Methods(http.MethodGet)
	router.HandleFunc("/register", user.Registration).Methods(http.MethodPost)
	router.HandleFunc("/login", user.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", user.Logout).Methods(http.MethodPost)

	router.HandleFunc("/chat", redirect)
	router.HandleFunc("/", redirect)

	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(authMiddleware)
	// PersonalChat
	protectedRouter.HandleFunc("/chat/find/{userName}", chats.FindChat)
	protectedRouter.HandleFunc("/chat/{userName}/history", chats.PersonalHistory)
	protectedRouter.HandleFunc("/chat/{userName}/ws", chats.PersonalChatWS)
	protectedRouter.HandleFunc("/chat/{userName}", chats.GetChat)

	// GlobalChat
	protectedRouter.HandleFunc("/globalChat", chats.GetGlobalChat)
	protectedRouter.HandleFunc("/globalChat/ws", chats.GlobalChatWS)
	protectedRouter.HandleFunc("/globalChat/history", chats.GlobalHistory)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Println("Start webMessage")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server:", err)
		}
	}()

	<-sigChan
	log.Println("Shutdown signal received, gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Error during server shutdown:", err)
	}

	if err := database.MongoClient.Disconnect(ctx); err != nil {
		log.Fatalln("MongoClient.Disconnect error:", err)
	}

	log.Println("Server stopped")
}
