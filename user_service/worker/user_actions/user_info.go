package user_actions

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func UserInfo(userInfoRequest requests.Request) requests.Response {
	req := userInfoRequest.(*requests.UserInfoRequest)
	row := userDataDB.QueryRow("SELECT username, email, bio, avatar_url FROM users WHERE id = $1", req.UserId)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return requests.NewErrorResponse(errors.UserDoesntExistError())
		}
		log.Println(row.Err())
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	r := new(requests.UserInfoResponse)
	if err := row.Scan(&r.Username, &r.Email, &r.Bio, &r.AvatarUrl); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return r
}
