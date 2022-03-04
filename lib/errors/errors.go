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
