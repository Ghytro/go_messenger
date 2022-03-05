package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
)

type Response interface {
	JsonBytes() []byte
	JsonString() string
	HTTPStatusCode() int
}

type ErrorResponse struct {
	Error          errors.Error `json:"error"`
	httpStatusCode int
}

func NewErrorResponse(err errors.Error) *ErrorResponse {
	return &ErrorResponse{err, err.HTTPStatusCode()}
}

func (er ErrorResponse) HTTPStatusCode() int {
	return er.httpStatusCode
}

func (er ErrorResponse) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(er)
	return jsonBytes
}

func (er ErrorResponse) JsonString() string {
	return string(er.JsonBytes())
}

type CreateUserResponse struct {
	httpStatusCode int
	Token          string `json:"token"`
}

func NewCreateUserResponse(token string) *CreateUserResponse {
	r := new(CreateUserResponse)
	r.httpStatusCode = http.StatusOK
	r.Token = token
	return r
}

func (cu CreateUserResponse) HTTPStatusCode() int {
	return cu.httpStatusCode
}

func (cu CreateUserResponse) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(cu)
	return jsonBytes
}

func (cu CreateUserResponse) JsonString() string {
	return string(cu.JsonBytes())
}

func SendResponse(w http.ResponseWriter, r Response) {
	http.Error(w, r.JsonString(), r.HTTPStatusCode())
}
