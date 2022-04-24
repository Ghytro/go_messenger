package message_actions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/message_service/worker/config"
)

func SendMessage(sendMessageRequest requests.Request) requests.Response {
	req := sendMessageRequest.(*requests.SendMessageRequest)

	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		log.Println(rdbGet.Err())
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()

	ctx := context.Background()
	tx, err := messageDataDB.BeginTx(ctx, nil)
	if err != nil {
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	defer tx.Rollback()
	// check if the user is in chat
	// returned members will be useful when creating notifications
	members := make([]int, 0)
	if err := tx.QueryRowContext(ctx,
		"SELECT members FROM chat_data WHERE id = $1 AND $2 = ANY(members)",
		req.ChatId,
		userId,
	).Scan(&members); err != nil {
		if err == sql.ErrNoRows {
			return requests.NewErrorResponse(errors.UnableToSendMessageError())
		}
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	// add message to special table with messages
	messageTableName := fmt.Sprintf("messages_%d", req.ChatId)
	if _, err := tx.ExecContext(ctx,
		fmt.Sprintf(
			`INSERT INTO %s
			(sender_id, message_text, attachments, parent_message)
			VALUES ($1, $2, $3, $4)`,
			messageTableName,
		),
		userId,
		req.MessageText,
		req.Attachments,
		req.ParentMessage,
	); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if err := tx.Commit(); err != nil {
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	// notify all users in this chat that the message was sent
	type Notification struct {
		Action      string                 `json:"type"`
		Description map[string]interface{} `json:"description"`
		Timestamp   int64                  `json:"timestamp"`
	}
	type JSONNotification struct {
		UserId       int          `json:"user_id"`
		Notification Notification `json:"notification"`
	}
	notifications := make([]JSONNotification, len(members))
	for i := range notifications {
		notifications[i] = JSONNotification{
			UserId: members[i],
			Notification: Notification{
				Action: "incoming_message",
				Description: map[string]interface{}{
					"chat_id":      req.ChatId,
					"from":         userId,
					"message_text": req.MessageText,
				},
				Timestamp: time.Now().Unix(),
			},
		}
	}
	jsonBytes, err := json.Marshal(notifications)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	response, err := http.Post(
		config.Config.NotificationServiceAddr,
		"application/json",
		strings.NewReader(string(jsonBytes)),
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Notification service response status: %d\n", response.StatusCode)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
