package chats

import (
	"fmt"
	"log"
	"net/http"
	"webMessenger/database"
	"webMessenger/user"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func PersonalChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("vars:", vars)

	user_name := vars["userName"]
	collection := database.GetCollection("users")
	var user user.User
	err := collection.FindOne(r.Context(), bson.D{{Key: "userName", Value: user_name}}).Decode(&user)
	if err != nil {
		log.Fatalln("Error PersonalChat (collection.FindOne()): ", err)
	} else {
		http.ServeFile(w, r, "./pages/UserNotFound.html")
		return
	}
}

