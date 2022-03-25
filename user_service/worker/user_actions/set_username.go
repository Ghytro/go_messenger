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
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	_, err := userDataDB.Exec(
		"UPDATE users SET username = $1 WHERE id = $2",
		req.Username,
		userId,
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
