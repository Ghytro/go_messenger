package message_actions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/message_service/worker/config"
)

func AddUser(addUserRequest requests.Request) requests.Response {
	req := addUserRequest.(*requests.AddUserRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	if userId == req.UserId {
		fmt.Println("here")
		return requests.NewErrorResponse(errors.UnableToInviteError())
	}
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	defer tx.Rollback()
	var chatId int
	err = tx.QueryRowContext(ctx, `
		SELECT id
		FROM chat_data
		WHERE id = $2 AND $1 != ALL(members) AND $1 != ALL(banned_users)
		`,
		req.UserId,
		req.ChatId,
	).Scan(&chatId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no rows", req.UserId, req.ChatId)
			return requests.NewErrorResponse(errors.UnableToInviteError())
		}
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}

	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE chat_data
		SET members = array_append(members, $2)
		WHERE id = $1
		`,
		chatId,
		req.UserId,
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	userServiceRequest := &requests.InviteUsersRequest{
		Token:        req.Token,
		ChatId:       req.ChatId,
		InvitedUsers: []int{req.UserId},
	}
	fmt.Println("made request to user service")
	userServiceResponse, err := httpClient.Post(
		config.Config.UserServiceAddr+"/invite_users",
		"application/json",
		strings.NewReader(userServiceRequest.JsonString()),
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	fmt.Println("made request to user service")
	switch sc := userServiceResponse.StatusCode; sc {
	case http.StatusInternalServerError:
		return requests.NewEmptyResponse(sc)
	case http.StatusBadRequest:
		fmt.Println("here3")
		responseBytes, err := io.ReadAll(userServiceResponse.Body)
		if err != nil {
			log.Println(err)
			return requests.NewEmptyResponse(http.StatusInternalServerError)
		}
		errResponse := requests.NewErrorResponse(errors.Error{}, sc)
		json.Unmarshal(responseBytes, errResponse)
		return errResponse
	}
	if err = tx.Commit(); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
