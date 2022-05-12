package user_actions

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

var allowedTokenSymbols = "0123456789qwertyuiopasdfghjklzxcvbnm"

const accessTokenLength = 50

func generateAccessToken() string {
	var tokenBuilder strings.Builder
	for i := 0; i < accessTokenLength; i++ {
		tokenBuilder.WriteByte(allowedTokenSymbols[rand.Int31n(int32(len(allowedTokenSymbols)))])
	}
	return tokenBuilder.String()
}

func checkMD5HashFormat(hash string) bool {
	if len(hash) != 32 {
		return false
	}
	for _, c := range hash {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

// too dumb for regexp
func checkUsernameFormat(username string) bool {
	if len(username) < 6 || len(username) > 20 {
		return false
	}
	containsLetters := false
	for _, c := range username {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			containsLetters = true
		} else if (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	if !containsLetters {
		return false
	}
	if strings.HasPrefix(username, "_") ||
		strings.HasSuffix(username, "_") ||
		strings.Contains(username, "__") {
		return false
	}
	return true
}

func checkEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && len(email) <= 320 // the length of email address is defined by international standart
}

func CreateUser(createUserRequest requests.Request) requests.Response {
	req := createUserRequest.(*requests.CreateUserRequest)
	if !checkUsernameFormat(req.Username) {
		return requests.NewErrorResponse(errors.IncorrectUsernameError())
	}
	if !checkEmailFormat(req.Email) {
		return requests.NewErrorResponse(errors.IncorrectEmailFormatError())
	}
	if !checkMD5HashFormat(req.PasswordMD5Hash) {
		return requests.NewErrorResponse(errors.IncorrectPasswordMD5HashError())
	}
	token := generateAccessToken()
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	var dbUsername, dbEmail, dbAccessToken string
	err = tx.QueryRowContext(
		ctx,
		`SELECT username, email, access_token
		FROM users WHERE username = $1 OR
		email = $2 OR
		access_token = $3`,
		req.Username,
		req.Email,
		token,
	).Scan(&dbUsername, &dbEmail, &dbAccessToken)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	for err == nil {
		if dbUsername == req.Username {
			return requests.NewErrorResponse(errors.UsernameAlreadyTakenError())
		}
		if dbEmail == req.Email {
			return requests.NewErrorResponse(errors.EmailAlreadyRegisteredError())
		}
		if dbAccessToken == token {
			token = generateAccessToken()
		}
		err = tx.QueryRowContext(
			ctx,
			`SELECT username, email, access_token
			FROM users WHERE username = $1 OR
			email = $2 OR
			access_token = $3`,
			req.Username,
			req.Email,
			token,
		).Scan(&dbUsername, &dbEmail, &dbAccessToken)
	}
	var newUserId int
	if err := tx.QueryRowContext(
		ctx,
		`INSERT INTO users
		(username, email, password_md5_hash, access_token)
		VALUES
		($1, $2, $3, $4)
		RETURNING id`,
		req.Username,
		req.Email,
		req.PasswordMD5Hash,
		token,
	).Scan(&newUserId); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if err = tx.Commit(); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	redisClient.Set(token, newUserId, 0)
	return requests.NewCreateUserResponse(token)
}
