package handler

import (
	"io"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/user_service/worker/user_actions"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	var response requests.Response
	switch r.URL.Path {
	case "/create_user":
		response = user_actions.CreateUser(requests.NewCreateUserRequest(bodyBytes))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPStatusCode())
	w.Write(response.JsonBytes())
}
