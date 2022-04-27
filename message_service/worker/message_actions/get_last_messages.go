package message_actions

import (
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

	if req.Amount < 0 || req.Offset < 0 {
		return requests.NewEmptyResponse(http.StatusBadRequest)
	}

	messagesTable := fmt.Sprint("messages_%d", req.ChatId)

	rows, err := messageDataDB.Query(
		fmt.Sprintf(
			"SELECT * FROM %s ORDER BY id DESC LIMIT $1 OFFSET $2",
			messagesTable,
		),
		req.Amount,
		req.Offset,
	)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == "" { // find out what the error code is
			return requests.NewErrorResponse(errors.InvalidChatIdError())
		}
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
		rows.Scan(
			&m.Id,
			&m.SenderId,
			&m.MessageText,
			pq.Array(&m.Attachments),
			&m.ParentMessage,
			&timestamp,
		)
		m.Timestamp = timestamp.Unix()
		response.Messages = append(response.Messages, m)
	}
	return response
}
