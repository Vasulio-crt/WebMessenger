package user

import (
	"encoding/json"
	"log"
	"net/http"
	"webMessenger/database"
	"webMessenger/user/utilities"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	UserName string `json:"userName" bson:"userName"`
	Password string `json:"password" bson:"password"`
}

func (u User) IsNull() bool {
	return u.UserName == "" || u.Password == ""
}

func GetRegistration(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/registration.html")
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/login.html")
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error server (json)", http.StatusInternalServerError)
		return
	}

	// Сохранения пользователя в БД
	if req.UserName == "" || req.Password == "" {
		http.Error(w, "fields are empty", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "error server (hash)", http.StatusInternalServerError)
		return
	}
	req.Password = hashedPassword

	collection := database.GetCollection("users")
	if _, err = collection.InsertOne(r.Context(), req); err != nil {
		http.Error(w, "error server (db)", http.StatusInternalServerError)
		return
	}

	// Сохранение сессии в БД
	err = create_session(w, r, req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect registration.html
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error server (json)", http.StatusInternalServerError)
		return
	}

	if req.UserName == "" || req.Password == "" {
		http.Error(w, "fields are empty", http.StatusBadRequest)
		return
	}

	var userDB User
	collection := database.GetCollection("users")
	err := collection.FindOne(r.Context(), bson.D{{Key: "userName", Value: req.UserName}}).Decode(&userDB)
	if err != nil {
		http.Error(w, "error server (db)", http.StatusInternalServerError)
		return
	}

	if !utilities.CheckPasswordHash(req.Password, userDB.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	err = create_session(w, r, req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect login.html
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Server cookie error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	err = delete_session(cookie)
	if err != nil {
		log.Println("Error logout:", err)
	}
}
