package main

import (
	"net/http"

	"github.com/Ghytro/go_messenger/user_service/interface/adapter"
	"github.com/Ghytro/go_messenger/user_service/interface/config"
)

func main() {
	for _, h := range config.Config.Handlers {
		http.HandleFunc("/"+h.Name, adapter.SendRequest)
	}
	http.ListenAndServe(":8082", nil)
}
