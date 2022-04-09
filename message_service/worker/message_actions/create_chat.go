package message_actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/message_service/worker/config"
	"github.com/lib/pq"
)

func CreateChat(createChatRequest requests.Request) requests.Response {
	req := createChatRequest.(*requests.CreateChatRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		log.Println(rdbGet.Err())
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	usersSlice := append([]int{userId}, req.Users.IntArray...)
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	row := tx.QueryRowContext(ctx,
		`INSERT INTO chat_data (name, avatar_url, admin_id, is_public, members)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		req.Name,
		req.AvatarUrl,
		userId,
		req.IsPublic,
		pq.Array(usersSlice),
	)
	if row.Err() != nil {
		fmt.Println("here")
		tx.Rollback()
		log.Println(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	var createdChatId int
	row.Scan(&createdChatId)
	messageTableName := fmt.Sprintf("messages_%d", createdChatId)
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY NOT NULL,
			sender_id INT NOT NULL,
			message_text TEXT,
			attachments VARCHAR(2048) [] NOT NULL DEFAULT '{}',
			send_timestamp TIMESTAMP NOT NULL,
			parent_message INT
		)`, messageTableName),
	)
	if err != nil {
		fmt.Println("here")
		tx.Rollback()
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}

	userServiceRequest := &requests.InviteUsersRequest{
		Token:        req.Token,
		InvitedUsers: usersSlice,
		ChatId:       createdChatId,
	}
	userServiceResponse, err := httpClient.Post(config.Config.UserServiceAddr+"/invite_users", "application/json", strings.NewReader(userServiceRequest.JsonString()))
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
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
