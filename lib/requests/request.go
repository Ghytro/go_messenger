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

type SetUsernameRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewSetUsernameRequest(jsonBytes []byte) *SetUsernameRequest {
	r := new(SetUsernameRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SetUsernameRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SetUsernameRequest) JsonString() string {
	return string(r.JsonBytes())
}

type SetPasswordRequest struct {
	Token              string `json:"token"`
	OldPasswordMD5Hash string `json:"old_password_md5_hash"`
	NewPasswordMD5Hash string `json:"new_password_md5_hash"`
}

func NewSetPasswordRequest(jsonBytes []byte) *SetPasswordRequest {
	r := new(SetPasswordRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SetPasswordRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SetPasswordRequest) JsonString() string {
	return string(r.JsonBytes())
}

type SetAvatarUrlRequest struct {
	Token     string `json:"token"`
	AvatarUrl string `json:"avatar_url"`
}

func NewSetAvatarUrlRequest(jsonBytes []byte) *SetAvatarUrlRequest {
	r := new(SetAvatarUrlRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SetAvatarUrlRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SetAvatarUrlRequest) JsonString() string {
	return string(r.JsonBytes())
}

type LogInRequest struct {
	Username        string `json:"username"`
	PasswordMD5Hash string `json:"password_md5_hash"`
}

func NewLogInRequest(jsonBytes []byte) *LogInRequest {
	r := new(LogInRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *LogInRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *LogInRequest) JsonString() string {
	return string(r.JsonBytes())
}
