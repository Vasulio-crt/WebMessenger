package user

import (
	"errors"
	"net/http"
	"webMessenger/database"
	"webMessenger/user/utilities"

	"go.mongodb.org/mongo-driver/bson"
)

type Session struct {
	SessionToken string `bson:"sessionToken"`
	UserName     string `bson:"userName"`
}

const collectionName = "sessions"

func create_session(w http.ResponseWriter, r *http.Request, userName string) error  {
	cookieValue := utilities.GenerateValueCookie()
	if cookieValue == "" {
		return errors.New("failed to generate cookie value")
	}

	collection := database.GetCollection(collectionName)
	session := Session{SessionToken: cookieValue, UserName: userName}
	_, err := collection.InsertOne(r.Context(), session)
	if err != nil {
		return errors.New("failed to save session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	})

	return nil
}

func Get_user_name(r *http.Request, cookie *http.Cookie) (string, error) {
	var session Session
	collection := database.GetCollection(collectionName)
	err := collection.FindOne(r.Context(), bson.D{{Key: "sessionToken", Value: cookie.Value}}).Decode(&session)
	if err != nil {
		return "", errors.New("failed find")
	}

	return session.UserName, nil
}