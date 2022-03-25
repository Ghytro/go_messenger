package handler

import (
	"io"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/message_service/worker/message_actions"
)

var HandlerMap = map[string]func(requests.Request) requests.Response{
	"/create_chat":         message_actions.CreateChat,
	"/set_chat_name":       message_actions.SetChatName,
	"/set_chat_avatar_url": message_actions.SetChatAvatarUrl,
	"/add_user":            message_actions.AddUser,
	"/kick_user":           message_actions.KickUser,
	"/ban_user":            message_actions.BanUser,
	"/chat_info":           message_actions.ChatInfo,
	"/send_message":        message_actions.SendMessage,
	"/delete_message":      message_actions.DeleteMessage,
}

var RequestGeneratorMap = map[string]func([]byte) requests.Request{
	"/create_chat": func(jb []byte) requests.Request {
		return requests.NewCreateChatRequest(jb)
	},
	"/set_chat_name": func(jb []byte) requests.Request {
		return requests.NewSetChatNameRequest(jb)
	},
	"/set_chat_avatar": func(jb []byte) requests.Request {
		return requests.NewSetChatAvatarUrlRequest(jb)
	},
	"/add_user": func(jb []byte) requests.Request {
		return requests.NewAddUserRequest(jb)
	},
	"/kick_user": func(jb []byte) requests.Request {
		return requests.NewKickUserRequest(jb)
	},
	"/ban_user": func(jb []byte) requests.Request {
		return requests.NewBanUserRequest(jb)
	},
	"/chat_info": func(jb []byte) requests.Request {
		return requests.NewChatInfoRequest(jb)
	},
	"/send_message": func(jb []byte) requests.Request {
		return requests.NewSendMessageRequest(jb)
	},
	"/delete_message": func(jb []byte) requests.Request {
		return requests.NewDeleteMessageRequest(jb)
	},
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	response := HandlerMap[r.URL.Path](RequestGeneratorMap[r.URL.Path](bodyBytes))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPStatusCode())
	w.Write(response.JsonBytes())
}
