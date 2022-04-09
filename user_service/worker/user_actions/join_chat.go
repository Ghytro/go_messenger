package user_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func JoinChat(joinChatRequest requests.Request) requests.Response {
	req := joinChatRequest.(*requests.JoinChatRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	_, err := userDataDB.Exec(
		"UPDATE users SET chats = array_append(chats, $1) WHERE id = $2",
		req.ChatId,
		userId,
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
