package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Ghytro/go_messenger/notification_service/notification"
)

func HandlePushNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	incomingNotifications := make(
		[]struct {
			UserId       int                       `json:"user_id"`
			Notification notification.Notification `json:"notification"`
		},
		0,
	)
	if err := json.Unmarshal(bodyBytes, &incomingNotifications); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, n := range incomingNotifications {
		notification.Push(
			n.UserId,
			n.Notification,
		)
	}
	w.WriteHeader(http.StatusOK)
}
