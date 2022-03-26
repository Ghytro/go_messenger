package user_actions

import (
	"context"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
)

func uniqueValues(slice []int) []int {
	type void struct{}
	uniqueMap := make(map[int]void)
	for _, v := range slice {
		uniqueMap[v] = void{}
	}
	result := make([]int, 0, len(uniqueMap))
	for k := range uniqueMap {
		result = append(result, k)
	}
	return result
}

func InviteUsers(inviteUsersRequest requests.Request) requests.Response {
	req := inviteUsersRequest.(*requests.InviteUsersRequest)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	userSlice := make([]int, 0)
	copy(userSlice, req.InvitedUsers)
	userSlice = append(userSlice, userId)
	userSlice = uniqueValues(userSlice)
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE user_chats
		SET chats = array_append(chats, $1)
		WHERE user_id IN $2 AND $1 NOT IN chats
		`,
		req.ChatId,
		userSlice,
	)

	if ra, _ := result.RowsAffected(); int(ra) != len(userSlice) {
		tx.Rollback()
		return requests.NewErrorResponse(errors.UnableToInviteError())
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
