package requests

import (
	"encoding/json"

	"github.com/Ghytro/go_messenger/lib/jsonhelpers"
)

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

type SetEmailRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewSetEmailRequest(jsonBytes []byte) *SetEmailRequest {
	r := new(SetEmailRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SetEmailRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SetEmailRequest) JsonString() string {
	return string(r.JsonBytes())
}

type UserInfoRequest struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

func NewUserInfoRequest(jsonBytes []byte) *UserInfoRequest {
	r := new(UserInfoRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *UserInfoRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *UserInfoRequest) JsonString() string {
	return string(r.JsonBytes())
}

type JoinChatRequest struct {
	Token  string `json:"token"`
	ChatId int    `json:"chat_id"`
}

func NewJoinChatRequest(jsonBytes []byte) *JoinChatRequest {
	r := new(JoinChatRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *JoinChatRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *JoinChatRequest) JsonString() string {
	return string(r.JsonBytes())
}

type InviteUserRequest struct {
	Token         string `json:"token"`
	ChatId        int    `json:"chat_id"`
	InvitedUserId int    `json:"invited_user_id"`
}

func NewInviteUserRequest(jsonBytes []byte) *InviteUserRequest {
	r := new(InviteUserRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *InviteUserRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *InviteUserRequest) JsonString() string {
	return string(r.JsonBytes())
}

type CreateChatRequest struct {
	Token     string                   `json:"token"`
	Name      string                   `json:"name"`
	AvatarUrl jsonhelpers.NullString   `json:"avatar_url"`
	Users     jsonhelpers.NullIntArray `json:"users"`
	IsPublic  bool                     `json:"is_public"`
}

func NewCreateChatRequest(jsonBytes []byte) *CreateChatRequest {
	r := new(CreateChatRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *CreateChatRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *CreateChatRequest) JsonString() string {
	return string(r.JsonBytes())
}

type SetChatNameRequest struct {
	Token    string `json:"token"`
	ChatId   int    `json:"chat_id"`
	ChatName string `json:"chat_name"`
}

func NewSetChatNameRequest(jsonBytes []byte) *SetChatNameRequest {
	r := new(SetChatNameRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SetChatNameRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SetChatNameRequest) JsonString() string {
	return string(r.JsonBytes())
}

type SetChatAvatarUrlRequest struct {
	Token     string `json:"token"`
	ChatId    int    `json:"chat_id"`
	AvatarUrl string `json:"chat_name"`
}

func NewSetChatAvatarUrlRequest(jsonBytes []byte) *SetChatAvatarUrlRequest {
	r := new(SetChatAvatarUrlRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SetChatAvatarUrlRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SetChatAvatarUrlRequest) JsonString() string {
	return string(r.JsonBytes())
}

type AddUserRequest struct {
	Token  string `json:"token"`
	ChatId int    `json:"chat_id"`
	UserId int    `json:"user_id"`
}

func NewAddUserRequest(jsonBytes []byte) *AddUserRequest {
	r := new(AddUserRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *AddUserRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *AddUserRequest) JsonString() string {
	return string(r.JsonBytes())
}

type KickUserRequest struct {
	Token  string `json:"token"`
	ChatId int    `json:"chat_id"`
	UserId int    `json:"user_id"`
}

func NewKickUserRequest(jsonBytes []byte) *KickUserRequest {
	r := new(KickUserRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *KickUserRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *KickUserRequest) JsonString() string {
	return string(r.JsonBytes())
}

type BanUserRequest struct {
	Token  string `json:"token"`
	ChatId int    `json:"chat_id"`
	UserId int    `json:"user_id"`
}

func NewBanUserRequest(jsonBytes []byte) *BanUserRequest {
	r := new(BanUserRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *BanUserRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *BanUserRequest) JsonString() string {
	return string(r.JsonBytes())
}

type ChatInfoRequest struct {
	Token   string `json:"token"`
	ChatIds []int  `json:"chat_ids"`
}

func NewChatInfoRequest(jsonBytes []byte) *ChatInfoRequest {
	r := new(ChatInfoRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *ChatInfoRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *ChatInfoRequest) JsonString() string {
	return string(r.JsonBytes())
}

type SendMessageRequest struct {
	Token         string                 `json:"token"`
	ChatId        int                    `json:"chat_id"`
	MessageText   jsonhelpers.NullString `json:"text"`
	Attachments   []string               `json:"attachments,omitempty"`
	ParentMessage jsonhelpers.NullInt    `json:"parent_message"`
}

func NewSendMessageRequest(jsonBytes []byte) *SendMessageRequest {
	r := new(SendMessageRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *SendMessageRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *SendMessageRequest) JsonString() string {
	return string(r.JsonBytes())
}

type DeleteMessageRequest struct {
	Token     string `json:"token"`
	ChatId    int    `json:"chat_id"`
	MessageId int    `json:"message_id"`
}

func NewDeleteMessageRequest(jsonBytes []byte) *DeleteMessageRequest {
	r := new(DeleteMessageRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *DeleteMessageRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *DeleteMessageRequest) JsonString() string {
	return string(r.JsonBytes())
}

type InviteUsersRequest struct {
	Token        string `json:"token"`
	InvitedUsers []int  `json:"invited_users"`
	ChatId       int    `json:"chat_id"`
}

func NewInviteUsersRequest(jsonBytes []byte) *InviteUsersRequest {
	r := new(InviteUsersRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *InviteUsersRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *InviteUsersRequest) JsonString() string {
	return string(r.JsonBytes())
}

type GetLastMessagesRequest struct {
	Token  string `json:"token"`
	ChatId int    `json:"chat_id"`
	Amount int    `json:"amount"`
	Offset int    `json:"offset"`
}

func NewGetLastMessagesRequest(jsonBytes []byte) *GetLastMessagesRequest {
	r := new(GetLastMessagesRequest)
	json.Unmarshal(jsonBytes, r)
	return r
}

func (r *GetLastMessagesRequest) JsonBytes() []byte {
	jsonBytes, _ := json.Marshal(*r)
	return jsonBytes
}

func (r *GetLastMessagesRequest) JsonString() string {
	return string(r.JsonBytes())
}
