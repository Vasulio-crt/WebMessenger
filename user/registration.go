package user

import (
	"encoding/json"
	"net/http"
	"webMessenger/database"
)

type User struct {
	UserName string `json:"userName" bson:"userName"`
	Password string `json:"password" bson:"password"`
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

	// Сохранение сессии в БД
	err := create_session(w, r, req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect на фронте
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
