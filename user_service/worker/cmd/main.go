package main

import (
	"fmt"
	"net/http"

	"github.com/Ghytro/go_messenger/user_service/worker/config"
	"github.com/Ghytro/go_messenger/user_service/worker/handler"
)

func main() {
	for _, h := range config.Config.ServedMethods {
		http.HandleFunc(h, handler.HandleRequest)
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
