package user_actions

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func validateUri(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil && len(uri) <= 2048
}

func SetAvatarUrl(setAvatarUrlRequest requests.Request) requests.Response {
	req := setAvatarUrlRequest.(*requests.SetAvatarUrlRequest)
	if !validateUri(req.AvatarUrl) {
		return requests.NewErrorResponse(errors.IncorrectUrlError())
	}
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	_, err := userDataDB.Exec(
		"UPDATE users SET avatar_url = $1 WHERE id = $2",
		req.AvatarUrl,
		userId,
	)
	if err != nil {
		log.Println(err) // debug
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
