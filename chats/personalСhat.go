package chats

import (
	"encoding/json"
	"log"
	"net/http"
	"webMessenger/database"
	"webMessenger/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

func GetChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	collection := database.GetCollection("users")
	var user user.User
	err := collection.FindOne(r.Context(), bson.D{{Key: "userName", Value: userName}}).Decode(&user)
	if err != nil {
		http.ServeFile(w, r, "./pages/UserNotFound.html")
		return
	}

	http.ServeFile(w, r, "./pages/chat.html")
}

func FindChat(w http.ResponseWriter, r *http.Request)  {
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

func PersonalHistory(w http.ResponseWriter, r *http.Request){

}

var hub = newHub()

func PersonalChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Error GlobalChat: ", err)
	}
	defer conn.Close()

	userName, err := user.Get_user_name(r, cookie)
	if err != nil {
		http.Error(w, "User name not found", http.StatusInternalServerError)
		return
	}

	client := &Client{
		Name: userName,
		Conn: conn,
	}
	hub.AddClient(client)
	defer hub.RemoveClient(client.Name)

	for {
		var msg MessageFromTo
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg.From = client.Name
		hub.SendPrivateMessage(msg)
	}
}
