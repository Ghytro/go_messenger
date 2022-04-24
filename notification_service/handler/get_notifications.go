package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/notification_service/notification"
	"github.com/go-redis/redis"
)

func HandleGetNotifications(w http.ResponseWriter, r *http.Request) {
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
	unmarshalled := struct {
		Token string `json:"token"`
	}{}
	if err := json.Unmarshal(tokenBytes, &unmarshalled); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := unmarshalled.Token
	userId, err := redisClient.Get(token).Int()
	if err != nil {
		if err == redis.Nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	n, err := notification.WaitForNotifications(userId)
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
