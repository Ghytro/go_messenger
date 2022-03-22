package user_actions

import (
	"log"
	"math/rand"
	"net/mail"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/lib/pq"
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
	_, err := userDataDB.Exec(
		`INSERT INTO users (username, email, password_md5_hash, access_token)
		VALUES ($1, $2, $3, $4)`,
		req.Username,
		req.Email,
		req.PasswordMD5Hash,
		token,
	)
	for err != nil {
		log.Println(err) // debug
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" { // unique constraint violation
			switch pqErr.Constraint {
			case "users_username_key":
				return requests.NewErrorResponse(errors.UsernameAlreadyTakenError())
			case "users_email_key":
				return requests.NewErrorResponse(errors.EmailAlreadyRegisteredError())
			case "users_access_token_key":
				token = generateAccessToken()
			}
		}
		_, err = userDataDB.Exec(`
		INSERT INTO users (username, email, password_md5_hash, access_token)
		VALUES ($1, $2, $3, $4)`,
			req.Username,
			req.Email,
			req.PasswordMD5Hash,
			token,
		)
	}
	redisClient.Set(token, req.Username, 0)
	return requests.NewCreateUserResponse(token)
}
