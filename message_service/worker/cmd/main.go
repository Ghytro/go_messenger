package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/message_service/worker/config"
	"github.com/Ghytro/go_messenger/message_service/worker/handler"
)

func main() {
	for handUrl := range handler.HandlerMap {
		http.HandleFunc(handUrl, handler.HandleRequest)
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
	log.Println(err)
}
