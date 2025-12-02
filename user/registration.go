package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webMessenger/database"
	"webMessenger/user/utilities"
)

type Session struct {
	SessionToken string `bson:"session_token"`
	UserName     string `bson:"user_name"`
}

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error server (json)", http.StatusInternalServerError)
		return
	}

	// Сохранения пользователя в БД
	if req.UserName == "" || req.Password == "" {
		http.Error(w, "the fields are empty", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("users")
	if _, err := collection.InsertOne(r.Context(), req); err != nil {
		http.Error(w, "error server (db)", http.StatusInternalServerError)
		return
	}
	
	cookieValue := utilities.GenerateValueCookie()
	if cookieValue == "" {
		http.Error(w, "Failed to generate cookie value", http.StatusInternalServerError)
		return
	}
	
	
	// Сохранение сессии в БД
	collection = database.GetCollection("sessions")
	session := Session{SessionToken: cookieValue, UserName: req.UserName}
	_, err := collection.InsertOne(r.Context(), session)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
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

	// TODO: сделать redirect
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
