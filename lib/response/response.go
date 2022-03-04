package response

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

type OKResponse struct {
}

func NewOKResponse() *OKResponse {
	return &OKResponse{}
}

func SendResponse(w http.ResponseWriter, r Response) {
	http.Error(w, r.JsonString(), r.HTTPStatusCode())
}
