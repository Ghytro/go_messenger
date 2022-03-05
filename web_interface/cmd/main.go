package main

import (
	"net/http"

	"github.com/Ghytro/go_messenger/web_interface/adapter"
	"github.com/Ghytro/go_messenger/web_interface/config"
)

func requestToUserService(w http.ResponseWriter, r *http.Request) {
	adapter.RequestToService(config.Config.UserServiceAddr, w, r)
}

func requestToChatService(w http.ResponseWriter, r *http.Request) {
	adapter.RequestToService(config.Config.ChatServiceAddr, w, r)
}

func main() {
	http.HandleFunc("/get_token", requestToUserService)
	http.HandleFunc("/revoke_token", requestToUserService)
	http.HandleFunc("/create_user", requestToUserService)
	http.HandleFunc("/whoami", requestToUserService)
	http.HandleFunc("/create_chat", requestToChatService)
	http.HandleFunc("/send_message", requestToChatService)
	http.ListenAndServe(":8080", nil)
}
