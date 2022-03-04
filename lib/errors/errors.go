package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		Code:           1,
		httpStatusCode: http.StatusBadRequest,
		Message:        "No access token provided to the API. Create if you don't have one or revoke the token.",
	}
}

func IncorrectHttpMethodError(expected string, got string) Error {
	return Error{
		Code:           2,
		httpStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("Incorrect http method. Expected: %s, but got: %s.", expected, got),
	}
}

func JsonValidationError() Error {
	return Error{
		Code:           3,
		httpStatusCode: http.StatusBadRequest,
		Message:        "Excepted JSON encoded data.",
	}
}

func MissingParameterError(parameter string) Error {
	return Error{
		Code:           4,
		httpStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("No required parameters in request: %s", parameter),
	}
}
