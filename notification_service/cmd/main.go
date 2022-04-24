package main

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/notification_service/handler"
)

func main() {
	clientRequestsMux, internalRequestsMux := http.NewServeMux(), http.NewServeMux()
	clientRequestsMux.HandleFunc("/", handler.HandleGetNotifications)
	internalRequestsMux.HandleFunc("/", handler.HandlePushNotifications)
	go func() {
		log.Println(http.ListenAndServe(":8079", clientRequestsMux))
	}()
	log.Println(http.ListenAndServe(":8078", internalRequestsMux))
}
