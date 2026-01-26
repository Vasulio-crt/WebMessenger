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

// cookieValue: userName
var allSessions = make(map[string]string, 0)

func create_session(w http.ResponseWriter, r *http.Request, userName string) error {
	cookieValue := utilities.GenerateValueCookie()
	if cookieValue == "" {
		return errors.New("failed to generate cookie value")
	}

	collection := database.GetCollection(collectionName)
	session := Session{SessionToken: cookieValue, UserName: userName}
	allSessions[cookieValue] = userName
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

func Init_sessions() error {
	collection := database.GetCollection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return errors.New("failed to find sessions")
	}
	defer cursor.Close(context.TODO())

	var sessions []Session
	if err = cursor.All(context.TODO(), &sessions); err != nil {
		return errors.New("failed to decode sessions")
	}

	for _, session := range sessions {
		allSessions[session.SessionToken] = session.UserName
	}

	return nil
}

func Get_user_name(cookie *http.Cookie) string {
	userName, ok := allSessions[cookie.Value]
	if ok {
		return userName
	}

	var session Session
	collection := database.GetCollection(collectionName)
	err := collection.FindOne(context.TODO(), bson.M{"sessionToken": cookie.Value}).Decode(&session)
	if err != nil {
		return ""
	}

	return session.UserName
}

func delete_session(cookie *http.Cookie) error {
	delete(allSessions, cookie.Value)
	collection := database.GetCollection(collectionName)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"sessionToken": cookie.Value})
	if err != nil {
		return errors.New("failed delete session")
	}
	return nil
}

// Вернет true Если сессия найдена, а иначе false
func Session_check(cookie *http.Cookie) bool {
	_, ok := allSessions[cookie.Value]
	if ok {
		return true
	}
	return false
}
