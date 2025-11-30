package user

import (
	"fmt"
	"net/http"
	"time"
	"webMessenger/database"
	"webMessenger/user/utilities"
)

type Session struct {
	SessionToken string    `bson:"session_token"`
	CreatedAt    time.Time `bson:"created_at"`
}

func Registration(w http.ResponseWriter, r *http.Request) {
	cookieValue := utilities.GenerateValueCookie()
	println("cookieValue", cookieValue) // debug
	if cookieValue == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to generate cookie value"))
		return
	}

	// Сохранение сессии в MongoDB
	collection := database.GetCollection("webMessenger", "sessions")
	session := Session{SessionToken: cookieValue, CreatedAt: time.Now()}
	_, err := collection.InsertOne(r.Context(), session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save session"))
		fmt.Println("Error inserting session to MongoDB:", err)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true, // Доступ только через HTTP, защита от XSS
		//Secure:   true, // Только HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, r *http.Request) {

}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
