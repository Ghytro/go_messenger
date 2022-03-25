package message_actions

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func CreateChat(createChatRequest requests.Request) requests.Response {
	req := createChatRequest.(*requests.CreateChatRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	var usersSlice []int
	if req.Users.Valid {
		usersSlice = req.Users.IntArray
	} else {
		usersSlice = make([]int, 0)
	}
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	row := tx.QueryRowContext(ctx, `
		INSERT INTO chat_data
		(name, avatar_url, admin_id, is_public, members)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING id`,
		req.Name,
		req.AvatarUrl,
		userId,
		req.IsPublic,
		usersSlice,
	)
	if row.Err() != nil {
		tx.Rollback()
		log.Println(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	var createdChatId int
	row.Scan(&createdChatId)
	messageTableName := fmt.Sprintf("messages_%d", createdChatId)
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE $1 (
			id SERIAL PRIMARY KEY NOT NULL,
			sender_id INT NOT NULL,
			message_text TEXT,
			attachments VARCHAR(2048) [] NOT NULL DEFAULT '{}',
			send_timestamp TIMESTAMP NOT NULL,
			parent_message INT
		)`,
		messageTableName,
	)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
