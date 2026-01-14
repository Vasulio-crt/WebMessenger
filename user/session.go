package user

import (
	"context"
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

var AllSessions = make(map[string]string)

func create_session(w http.ResponseWriter, r *http.Request, userName string) error  {
	cookieValue := utilities.GenerateValueCookie()
	if cookieValue == "" {
		return errors.New("failed to generate cookie value")
	}

	collection := database.GetCollection(collectionName)
	session := Session{SessionToken: cookieValue, UserName: userName}
	AllSessions[cookieValue] = userName
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
	userName, ok := AllSessions[cookie.Value]
	if ok {
		return userName, nil
	}
	
	var session Session
	collection := database.GetCollection(collectionName)
	err := collection.FindOne(r.Context(), bson.D{{Key: "sessionToken", Value: cookie.Value}}).Decode(&session)
	if err != nil {
		return "", errors.New("failed find")
	}

	return session.UserName, nil
}

func delete_session(cookie *http.Cookie) error {
	delete(AllSessions, cookie.Value)
	collection := database.GetCollection(collectionName)
	_, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "sessionToken", Value: cookie.Value}})
	if err != nil {
		return errors.New("failed delete session")
	}
	return nil
}