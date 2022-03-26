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
	result, err := userDataDb.Exec(
		"UPDATE user_chats SET chats = array_append(chats, $1) WHERE user_id = $2 AND $1 NOT IN chats",
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
