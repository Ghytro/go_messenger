package main

import (
	"encoding/json"
	"fmt"
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
	rdbGet := redisClient.Get(token)
	if rdbGet.Err() != nil {
		log.Println(rdbGet.Err())
		return false
	}
	id, _ := rdbGet.Int()
	return id != 0
}

func handleGetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	unmarshalled := make(map[string]interface{})
	if err := json.Unmarshal(tokenBytes, &unmarshalled); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := unmarshalled["token"].(string)
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
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var incomingNotification struct {
		Token        string                    `json:"token"`
		Notification notification.Notification `json:"notification"`
	}
	if err := json.Unmarshal(bodyBytes, &incomingNotification); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("here")
	notification.Push(
		incomingNotification.Token,
		incomingNotification.Notification,
	)
	fmt.Println("here")
	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println(config.Config.RedisTokenValidationAddr)
	fmt.Println(redisClient.Ping().String())
	clientRequestsMux, internalRequestsMux := http.NewServeMux(), http.NewServeMux()
	clientRequestsMux.HandleFunc("/", handleGetNotifications)
	internalRequestsMux.HandleFunc("/", handlePushNotification)
	go func() {
		log.Println(http.ListenAndServe(":8079", clientRequestsMux))
	}()
	log.Println(http.ListenAndServe(":8078", internalRequestsMux))
}
