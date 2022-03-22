package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/jsonhelpers"
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

type EmptyResponse struct {
	httpStatusCode int
}

func NewEmptyResponse(statusCode int) *EmptyResponse {
	return &EmptyResponse{statusCode}
}

func (er EmptyResponse) HTTPStatusCode() int {
	return er.httpStatusCode
}

func (er EmptyResponse) JsonBytes() []byte {
	return []byte(er.JsonString())
}

func (er EmptyResponse) JsonString() string {
	return ""
}

type CreateUserResponse struct {
	Token string `json:"token"`
}

func NewCreateUserResponse(token string) *CreateUserResponse {
	r := new(CreateUserResponse)
	r.Token = token
	return r
}

func (cu CreateUserResponse) HTTPStatusCode() int {
	return http.StatusOK
}

func (cu CreateUserResponse) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(cu)
	return jsonBytes
}

func (cu CreateUserResponse) JsonString() string {
	return string(cu.JsonBytes())
}

type LogInResponse struct {
	Token string `json:"token"`
}

func NewLogInResponse(token string) *LogInResponse {
	return &LogInResponse{token}
}

func (lr LogInResponse) HTTPStatusCode() int {
	return http.StatusOK
}

func (lr LogInResponse) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(lr)
	return jsonBytes
}

func (lr LogInResponse) JsonString() string {
	return string(lr.JsonBytes())
}

type UserInfoResponse struct {
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	AvatarUrl jsonhelpers.NullString `json:"avatar_url"`
	Bio       jsonhelpers.NullString `json:"bio"`
}

func (ur UserInfoResponse) HTTPStatusCode() int {
	return http.StatusOK
}

func (ur UserInfoResponse) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(ur)
	return jsonBytes
}

func (ur UserInfoResponse) JsonString() string {
	return string(ur.JsonBytes())
}

func SendResponse(w http.ResponseWriter, r Response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.HTTPStatusCode())
	w.Write(r.JsonBytes())
}
