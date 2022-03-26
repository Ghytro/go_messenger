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
	"github.com/Ghytro/go_messenger/message_service/config"
)

func AddUser(addUserRequest requests.Request) requests.Response {
	req := addUserRequest.(*requests.AddUserRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE chat_data
		SET members = array_append(members, $1)
		WHERE id = $2 AND $1 NOT IN members AND $1 NOT IN banned_users
		`,
		req.UserId,
		req.ChatId,
	)
	if ra, _ := result.RowsAffected(); ra == 0 {
		return requests.NewErrorResponse(errors.UnableToInviteError)
	}
	userServiceRequest := &requests.AddUserRequest{
		Token:  req.Token,
		ChatId: req.ChatId,
		UserId: req.UserId,
	}
	userServiceResponse, err := httpClient.Post(config.Config.UserServiceAddr+"/add_user", "application/json", strings.NewReader(userServiceRequest.JsonString()))
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	switch sc := userServiceResponse.StatusCode; sc {
	case http.StatusInternalServerError:
		tx.Rollback()
		return requests.NewEmptyResponse(sc)
	case http.StatusBadRequest:
		responseBytes, err := io.ReadAll(userServiceResponse.Body)
		if err != nil {
			log.Println(err)
			return requests.NewEmptyResponse(http.StatusInternalServerError)
		}
		errResponse := new(requests.ErrorResponse)
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
