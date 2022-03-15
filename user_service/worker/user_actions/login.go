package user_actions

import (
	"database/sql"
	"log"
	"net/http"
)

func Login(req *requests.LogInRequest) requests.Reponse {
	row := userDataDB.QueryRow("SELECT token, password_md5_hash FROM users WHERE username = $1", req.Username)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return requests.NewErrorResponse(errors.NoUserFound())
		}
		log.Fatal(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	var token, passwordMD5Hash string
	if err = row.Scan(&token, &password_md5_hash); err != nil {
		log.Fatal(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if passwordMD5Hash != req.PasswordMD5Hash {
		return requests.NewErrorResponse(errors.WrongPasswordError())
	}
	return requests.NewLogInResponse(token)

