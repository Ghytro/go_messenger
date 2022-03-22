package user_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func SetEmail(setEmailRequest requests.Request) requests.Response {
	req := setEmailRequest.(*requests.SetEmailRequest)
	if !checkEmailFormat(req.Email) {
		return requests.NewErrorResponse(errors.IncorrectEmailFormatError())
	}
	_, err := userDataDB.Exec(
		"UPDATE users SET email = $1 WHERE access_token = $2",
		req.Email,
		req.Token,
	)
	if err != nil {
		log.Println(err) // debug
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
