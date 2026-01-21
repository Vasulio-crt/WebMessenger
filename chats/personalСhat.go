package chats

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"webMessenger/database"
	"webMessenger/user"
	"webMessenger/user/utilities"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || !user.Session_check(cookie) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	userName := vars["userName"]
	collection := database.GetCollection("users")
	var user user.User
	err = collection.FindOne(r.Context(), bson.D{{Key: "userName", Value: userName}}).Decode(&user)
	if err != nil {
		http.ServeFile(w, r, "./pages/UserNotFound.html")
		return
	}

	http.ServeFile(w, r, "./pages/chat.html")
}

func FindChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	collection := database.GetCollection("users")
	var user user.User
	err := collection.FindOne(r.Context(), bson.D{{Key: "userName", Value: userName}}).Decode(&user)
	response := make(map[string]bool, 1)

	if err == nil && !user.IsNull() {
		response["found"] = true
	} else {
		response["found"] = false
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PersonalHistory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || !user.Session_check(cookie) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	users := []string{mux.Vars(r)["userName"]}
	user1 := user.Get_user_name(cookie)
	if user1 == "" {
		http.Error(w, "User name not found", http.StatusInternalServerError)
		return
	}
	users = append(users, user1)
	sort.Strings(users)
	hashChat := utilities.HashString(users[0] + users[1])

	collection := database.GetCollectionHistory(hashChat)
	cursor, err := collection.Find(r.Context(), bson.D{}, options.Find().SetLimit(60))
	if err != nil {
		http.Error(w, "Failed to retrieve message history", http.StatusInternalServerError)
		log.Println("Error finding messages in MongoDB:", err)
		return
	}
	defer cursor.Close(r.Context())

	var chatHistory []Message
	if err = cursor.All(r.Context(), &chatHistory); err != nil {
		http.Error(w, "Failed to decode message history", http.StatusInternalServerError)
		log.Println("Error decoding messages from MongoDB:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatHistory)
}

var hub = newHub()

func PersonalChatWS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie || !user.Session_check(cookie) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Error PersonalChat: ", err)
	}
	defer conn.Close()

	userName := user.Get_user_name(cookie)
	if userName == "" {
		http.Error(w, "User name not found", http.StatusInternalServerError)
		return
	}

	hub.AddClient(userName, conn)
	defer hub.RemoveClient(userName)

	for {
		var msg MessageFromTo
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg.From = userName
		hub.SendPrivateMessage(msg)
	}
}
