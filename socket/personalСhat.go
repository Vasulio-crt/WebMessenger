package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func PersonalChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("vars:", vars)

	user_name := vars["user_name"]
	w.Write([]byte(fmt.Sprintf("This is a personal chat for user: %s", user_name)))
}