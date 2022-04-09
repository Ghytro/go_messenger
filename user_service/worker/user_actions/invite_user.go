package user_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func InviteUser(inviteUserRequest requests.Request) requests.Response {
	req := inviteUserRequest.(*requests.InviteUserRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	if userId == req.InvitedUserId {
		return requests.NewErrorResponse(errors.UnableToInviteError())
	}
	result, err := userDataDB.Exec(
		"UPDATE users SET chats = array_append(chats, $1) WHERE id = $2 AND $1 NOT IN chats",
		req.ChatId,
		req.InvitedUserId,
	)
	if ra, _ := result.RowsAffected(); ra == 0 {
		return requests.NewErrorResponse(errors.UnableToInviteError())
	}
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
