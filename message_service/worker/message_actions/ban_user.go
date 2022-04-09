package message_actions

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/message_service/worker/config"
)

func BanUser(banUserRequest requests.Request) requests.Response {
	req := banUserRequest.(*requests.BanUserRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	if userId == req.UserId {
		return requests.NewErrorResponse(errors.UnableToBanError())
	}
	ctx := context.Background()
	tx, err := messageDataDB.BeginTx(ctx, nil)
	// todo: add to ban list and remove chat in user service
	result, err := tx.ExecContext(ctx, `
		UPDATE chat_data
		SET
			banned_users = array_append(banned_users, $1),
			members = ARRAY(SELECT unnest(members) EXCEPT SELECT $1)
		WHERE
			id = $2 AND
			$1 IN members AND
			admin_id = $3`,
		req.UserId,
		req.ChatId,
		userId,
	)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if ra, _ := result.RowsAffected(); ra == 0 {
		tx.Rollback()
		return requests.NewErrorResponse(errors.UnableToBanError())
	}
	userServiceResponse, err := httpClient.Post(config.Config.UserServiceAddr+"/remove_from_chat", "application/json", strings.NewReader(req.JsonString()))
	switch sc := userServiceResponse.StatusCode; sc {
	case http.StatusInternalServerError:
		return requests.NewEmptyResponse(sc)
	case http.StatusBadRequest:
		errResponse := new(requests.ErrorResponse)
		responseBytes, err := io.ReadAll(userServiceResponse.Body)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return requests.NewEmptyResponse(http.StatusInternalServerError)
		}
		json.Unmarshal(responseBytes, errResponse)
		return errResponse
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
