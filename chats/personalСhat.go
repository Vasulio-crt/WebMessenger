package chats

import (
	"fmt"
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
	fmt.Println("vars:", vars)

	user_name := vars["userName"]
	collection := database.GetCollection("users")
	var user user.User
	err := collection.FindOne(r.Context(), bson.D{{Key: "userName", Value: user_name}}).Decode(&user)
	fmt.Println("user:", user, user.IsNull())
	if err != nil || user.IsNull() {
		http.ServeFile(w, r, "./pages/UserNotFound.html")
		return
	}

	// Если пользователь найден, здесь должна быть логика для отображения чата.
	// Например, http.ServeFile(w, r, "./pages/personalChat.html")
}

var hub = newHub()

func PersonalChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
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
		var msg Message
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
