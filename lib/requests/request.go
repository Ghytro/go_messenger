package requests

import "encoding/json"

type Request interface {
	JsonBytes() []byte
	JsonString() string
}

type CreateUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	PasswordMD5Hash string `json:"password_md5_hash"`
}

func NewCreateUserRequest(jsonBytes []byte) *CreateUserRequest {
	r := new(CreateUserRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *CreateUserRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *CreateUserRequest) JsonString() string {
	return string(r.JsonBytes())
}
