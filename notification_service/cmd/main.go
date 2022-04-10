package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/notification_service/config"
	"github.com/Ghytro/go_messenger/notification_service/notification"
	"github.com/go-redis/redis"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     config.Config.RedisTokenValidationAddr,
	Password: "",
	DB:       0,
})

func verifyToken(token string) bool {
	userId, err := redisClient.Get(token).Int()
	return err == nil && userId != 0
}

func handleGetNotifications(w http.ResponseWriter, r *http.Request) {
	tokenBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := string(tokenBytes)
	if !verifyToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	n, err := notification.WaitForNotifications(token)
	if err != nil { // timeout
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}
	responseBytes, err := json.Marshal(notification.Notifications{List: n})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(responseBytes)
}

func handlePushNotification(w http.ResponseWriter, r *http.Request) {

}

func main() {
	clientRequestsMux, internalRequestsMux := http.NewServeMux(), http.NewServeMux()
	clientRequestsMux.HandleFunc("/", handleGetNotifications)
	internalRequestsMux.HandleFunc("/", handlePushNotification)
	go func() {
		log.Println(http.ListenAndServe(":8079", clientRequestsMux))
	}()
	log.Println(http.ListenAndServe(":8078", internalRequestsMux))
}
