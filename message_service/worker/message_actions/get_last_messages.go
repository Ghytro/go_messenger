package message_actions

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/lib/pq"
)

func GetLastMessages(getLastMessagesRequest requests.Request) requests.Response {
	req := getLastMessagesRequest.(*requests.GetLastMessagesRequest)

	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		log.Println(rdbGet.Err())
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	fmt.Println(userId)

	if req.Amount < 0 || req.Offset < 0 {
		return requests.NewEmptyResponse(http.StatusBadRequest)
	}
	if req.Amount == 0 {
		fmt.Println("amount 0")
		return &requests.GetLastMessagesResponse{
			ChatId:   req.ChatId,
			Messages: make([]requests.Message, 0),
		}
	}

	ctx := context.Background()
	tx, err := messageDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	if err := messageDataDB.QueryRowContext(
		ctx,
		"SELECT id FROM chat_data WHERE id = $1 AND (is_public = TRUE OR $2 = ANY(members))",
		req.ChatId,
		userId,
	).Err(); err != nil {
		if err == sql.ErrNoRows {
			return requests.NewErrorResponse(errors.UnableToGetMessagesError())
		}
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}

	messagesTable := fmt.Sprintf("messages_%d", req.ChatId)
	rows, err := tx.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT * FROM %s ORDER BY id DESC LIMIT $1 OFFSET $2",
			messagesTable,
		),
		req.Amount,
		req.Offset,
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	response := &requests.GetLastMessagesResponse{
		ChatId:   req.ChatId,
		Messages: make([]requests.Message, 0),
	}
	for rows.Next() {
		var m requests.Message
		var timestamp time.Time
		if err := rows.Scan(
			&m.Id,
			&m.SenderId,
			&m.MessageText,
			pq.Array(&m.Attachments),
			&timestamp,
			&m.ParentMessage,
		); err != nil {
			log.Println(err)
			return requests.NewEmptyResponse(http.StatusInternalServerError)
		}
		m.Timestamp = timestamp.Unix()
		response.Messages = append(response.Messages, m)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return response
}
