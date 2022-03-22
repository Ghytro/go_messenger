package user_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/lib/pq"
)

func SetUsername(setUsernameRequest requests.Request) requests.Response {
	req := setUsernameRequest.(*requests.SetUsernameRequest)
	if !checkUsernameFormat(req.Username) {
		return requests.NewErrorResponse(errors.IncorrectUsernameError())
	}
	_, err := userDataDB.Exec(
		"UPDATE users SET username = $1 WHERE access_token = $2",
		req.Username,
		req.Token,
	)
	if err != nil {
		log.Println(err) // debug
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" { // unique constraint violation
			switch pqErr.Constraint {
			case "users_username_key":
				return requests.NewErrorResponse(errors.UsernameAlreadyTakenError())
			}
		}
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
