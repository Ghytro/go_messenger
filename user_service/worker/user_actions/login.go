package user_actions

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func LogIn(logInRequest requests.Request) requests.Response {
	req := logInRequest.(*requests.LogInRequest)
	row := userDataDB.QueryRow("SELECT access_token, password_md5_hash FROM users WHERE username = $1", req.Username)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return requests.NewErrorResponse(errors.UserDoesntExistError())
		}
		log.Println(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	var token, passwordMD5Hash string
	if err := row.Scan(&token, &passwordMD5Hash); err != nil {
		log.Fatal(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if passwordMD5Hash != req.PasswordMD5Hash {
		return requests.NewErrorResponse(errors.WrongPasswordError())
	}
	return requests.NewLogInResponse(token)
}
