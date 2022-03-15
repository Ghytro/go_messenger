package user_actions

import (
	"context"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func SetPassword(req *requests.SetPasswordRequest) requests.Response {
	if !checkMD5HashFormat(req.OldPasswordMD5Hash) || !checkMD5HashFormat(req.NewPasswordMD5Hash) {
		return requests.NewErrorResponse(errors.IncorrectPasswordMD5HashError())
	}
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	row := tx.QueryRowContext("SELECT password_md5_hash FROM users WHERE token = $1", req.Token)
	if row.Err() != nil {
		log.Fatal(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	var oldPassMD5Hash string
	err = row.Scan(&oldPassMD5Hash)
	if err != nil {
		log.Fatal(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if oldPassMD5Hash != req.OldPasswordMD5Hash {
		return requests.NewErrorResponse(errors.OldPasswordsDoesntMatchError())
	}
	newToken := generateAccessToken()
	_, err := tx.ExecContext("UPDATE users SET password_md5_hash = $1, token = $2 WHERE token = $3", req.NewPasswordMD5Hash, newToken, req.Token)
	while err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_access_token_key":
				token = generateAccessToken()
			}
		}
		_, err := tx.ExecContext("UPDATE users SET password_md5_hash = $1, token = $2 WHERE token = $3", req.NewPasswordMD5Hash, newToken, req.Token)
	}
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
}
