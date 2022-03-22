package main

import (
	"net/http"

	"github.com/Ghytro/go_messenger/web_interface/adapter"
	"github.com/Ghytro/go_messenger/web_interface/config"
)

func requestToUserService(w http.ResponseWriter, r *http.Request) {
	adapter.RequestToService(config.Config.UserServiceAddr, w, r)
}

func requestToMessageService(w http.ResponseWriter, r *http.Request) {
	adapter.RequestToService(config.Config.MessageServiceAddr, w, r)
}

func main() {
	for _, m := range config.Config.UserServiceMethods {
		http.HandleFunc(m, requestToUserService)
	}
	for _, m := range config.Config.MessageServiceMethods {
		http.HandleFunc(m, requestToMessageService)
	}
	http.ListenAndServe(":8080", nil)
}
