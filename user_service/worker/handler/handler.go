package handler

import (
	"io"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/user_service/worker/user_actions"
)

var HandlerMap = map[string]func(requests.Request) requests.Response{
	"/create_user":    user_actions.CreateUser,
	"/set_username":   user_actions.SetUsername,
	"/set_password":   user_actions.SetPassword,
	"/set_avatar_url": user_actions.SetAvatarUrl,
	"/user_info":      user_actions.UserInfo,
	"/set_email":      user_actions.SetEmail,
	"/login":          user_actions.LogIn,
}

var RequestGeneratorMap = map[string]func([]byte) requests.Request{
	"/create_user": func(jb []byte) requests.Request {
		return requests.NewCreateUserRequest(jb)
	},
	"/set_username": func(jb []byte) requests.Request {
		return requests.NewSetUsernameRequest(jb)
	},
	"/set_password": func(jb []byte) requests.Request {
		return requests.NewSetPasswordRequest(jb)
	},
	"/set_avatar_url": func(jb []byte) requests.Request {
		return requests.NewSetAvatarUrlRequest(jb)
	},
	"/user_info": func(jb []byte) requests.Request {
		return requests.NewUserInfoRequest(jb)
	},
	"/set_email": func(jb []byte) requests.Request {
		return requests.NewSetEmailRequest(jb)
	},
	"/login": func(jb []byte) requests.Request {
		return requests.NewLogInRequest(jb)
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
