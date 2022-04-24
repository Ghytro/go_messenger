package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	NoAccessTokenErrCode = iota + 1
	InvalidAccessTokenErrCode
	IncorrectHttpMethodErrCode
	JsonValidationErrorCode
	MissingParameterErrorCode
	UsernameAlreadyTakenErrorCode
	IncorrectUsernameErrorCode
	EmailAlreadyRegisteredErrorCode
	IncorrectEmailFormatErrorCode
	IncorrectPasswordMD5HashErrorCode
	OldPasswordsDoesntMatchErrorCode
	UserDoesntExistErrorCode
	WrongPasswordErrorCode
	IncorrectUrlErrorCode
	UnableToInviteErrorCode
	UnableToBanErrorCode
	UnableToSendMessageErrorCode
)

type Error struct {
	Code           int `json:"code"`
	httpStatusCode int
	Message        string `json:"message"`
}

func (e Error) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(e)
	return jsonBytes
}

func (e Error) JsonString() string {
	return string(e.JsonBytes())
}

func (e Error) HTTPStatusCode() int {
	return e.httpStatusCode
}

func NoAccessTokenError() Error {
	return Error{
		Code:           NoAccessTokenErrCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "No access token provided to the API. Create if you don't have one or revoke the token.",
	}
}

func InvalidAccessTokenError() Error {
	return Error{
		Code:           InvalidAccessTokenErrCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Invalid api token. Check the sent token or try to revoke the token.",
	}
}

func IncorrectHttpMethodError(expected string, got string) Error {
	return Error{
		Code:           IncorrectHttpMethodErrCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("Incorrect http method. Expected: %s, but got: %s.", expected, got),
	}
}

func JsonValidationError() Error {
	return Error{
		Code:           JsonValidationErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Excepted JSON encoded data.",
	}
}

func MissingParameterError(parameter string) Error {
	return Error{
		Code:           MissingParameterErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("No required parameters in request: %s", parameter),
	}
}

func UsernameAlreadyTakenError() Error {
	return Error{
		Code:           UsernameAlreadyTakenErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Username is already taken, try another one.",
	}
}

func IncorrectUsernameError() Error {
	return Error{
		Code:           UsernameAlreadyTakenErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Your username doesn't require given limitations: username should have length from 6 to 20 symbols, should not start or end with underscores, have double underscores or contain symbols except english characters, digits or underscores.",
	}
}

func EmailAlreadyRegisteredError() Error {
	return Error{
		Code:           EmailAlreadyRegisteredErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "An account is already registered to this email.",
	}
}

func IncorrectEmailFormatError() Error {
	return Error{
		Code:           IncorrectEmailFormatErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Your email has incorrect format",
	}
}

func IncorrectPasswordMD5HashError() Error {
	return Error{
		Code:           IncorrectPasswordMD5HashErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Incorrect format of MD5 hash.",
	}
}

func OldPasswordsDoesntMatchError() Error {
	return Error{
		Code:           OldPasswordsDoesntMatchErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "MD5 hash of old password in request doesnt match to the actual one.",
	}
}

func UserDoesntExistError() Error {
	return Error{
		Code:           UserDoesntExistErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "No such user exists. Check the correctness of data in query.",
	}
}

func WrongPasswordError() Error {
	return Error{
		Code:           WrongPasswordErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Wrong password.",
	}
}

func IncorrectUrlError() Error {
	return Error{
		Code:           IncorrectUrlErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Incorrect format of url",
	}
}

func UnableToInviteError() Error {
	return Error{
		Code:           UnableToInviteErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Unable to invite some of the users specified in request. Check if they are already in chat and all user ids are valid.",
	}
}

func UnableToBanError() Error {
	return Error{
		Code:           UnableToBanErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Unable to ban the specified user. Check the correctness of sent request.",
	}
}

func UnableToSendMessageError() Error {
	return Error{
		Code:           UnableToSendMessageErrorCode,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Unable to send message to this chat",
	}
}
