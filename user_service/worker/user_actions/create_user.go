package user_actions

import (
	"log"
	"math/rand"
	"net/mail"
	"regexp"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/lib/sqlhelpers"
	"github.com/lib/pq"
)

var allowedTokenSymbols = "0123456789qwertyuiopasdfghjklzxcvbnm"

const accessTokenLength = 50

func generateAccessToken() string {
	var tokenBuilder strings.Builder
	for i := 0; i < accessTokenLength; i++ {
		tokenBuilder.WriteByte(allowedTokenSymbols[rand.Int31n(accessTokenLength)])
	}
	return tokenBuilder.String()
}

func checkMD5HashFormat(hash string) bool {
	result, err := regexp.Match("/^[a-f0-9]{32}$/i", []byte(hash))
	if err != nil {
		log.Fatal(err) // debug
	}
	return result
}

func checkUsernameFormat(username string) bool {
	result, err := regexp.Match("^(?=[a-zA-Z0-9._]{6,20}$)(?!.*[_.]{2})[^_.].*[^_.]$", []byte(username))
	if err != nil {
		log.Fatal(err) // debug
	}
	return result
}

func checkEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && len(email) <= 320 // the length of email address is defined by international standart
}

func CreateUser(req *requests.CreateUserRequest) requests.Response {
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
	_, err := sqlhelpers.RunQueryFromFile(userDataDB, "sql/create_user.sql", req.Username, req.Email, req.PasswordMD5Hash, token)
	for err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" { // unique constraint violation
			switch pqErr.Constraint {
			case "users_pkey":
				return requests.NewErrorResponse(errors.UsernameAlreadyTakenError())
			case "users_email_key":
				return requests.NewErrorResponse(errors.EmailAlreadyRegisteredError())
			case "users_access_token_key":
				token = generateAccessToken()
			}
		}
		_, err = sqlhelpers.RunQueryFromFile(userDataDB, "sql/create_user.sql", req.Username, req.Email, req.PasswordMD5Hash, token)
	}
	redisClient.Set(token, req.Username, 0)
	return requests.NewCreateUserResponse(token)
}
