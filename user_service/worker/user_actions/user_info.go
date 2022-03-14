package user_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/requests"
)

func UserInfo(req *requests.UserInfoRequest) requests.Response {
	row := userDataDB.QueryRow("SELECT username, email, bio, avatar_url FROM users WHERE username = $1", req.Username)
	if row.Err() != nil {
		log.Fatal(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	r := new(requests.UserInfoResponse)
	if err := row.Scan(&r.Username, &r.Email, &r.Bio, &r.AvatarUrl); err != nil {
		log.Fatal(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return r
}
