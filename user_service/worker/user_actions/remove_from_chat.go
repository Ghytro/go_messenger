package user_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func RemoveFromChat(banUserRequest requests.Request) requests.Response {
	req := banUserRequest.(*requests.BanUserRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		log.Println(rdbGet.Err())
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	result, err := userDataDB.Exec(`
		UPDATE user_chats
		SET chats = ARRAY(SELECT unnest(chats) EXCEPT SELECT $1)
		WHERE user_id = $2`,
		req.ChatId,
		req.UserId,
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if ra, _ := result.RowsAffected(); ra == 0 {
		return requests.NewErrorResponse(errors.UnableToBanError())
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
