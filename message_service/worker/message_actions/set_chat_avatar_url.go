package message_actions

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func SetChatAvatarUrl(setChatAvatarUrlRequest requests.Request) requests.Response {
	req := setChatAvatarUrlRequest.(*requests.SetChatAvatarUrlRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	if userId == 0 {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	ctx := context.Background()
	tx, err := messageDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx, `
		SELECT id
		FROM chat_data
		WHERE id = $1 AND admin_id = $2
		`,
		req.ChatId,
		userId,
	)
	if err := row.Err(); err != nil {
		if err == sql.ErrNoRows {
			return requests.NewEmptyResponse(http.StatusForbidden)
		}
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE chat_data
		SET avatar_url = $1
		WHERE id = $2
		`,
		req.AvatarUrl,
		req.ChatId,
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
