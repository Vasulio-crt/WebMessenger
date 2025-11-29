package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func PersonalChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("vars:", vars)
}