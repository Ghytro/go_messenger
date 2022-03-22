package main

import (
	"fmt"
	"net/http"

	"github.com/Ghytro/go_messenger/user_service/worker/config"
	"github.com/Ghytro/go_messenger/user_service/worker/handler"
)

func main() {
	for handUrl := range handler.HandlerMap {
		http.HandleFunc(handUrl, handler.HandleRequest)
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
