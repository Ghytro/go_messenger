package user_actions

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/lib/pq"
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
	fmt.Println(*req)
	rdbGet := redisClient.Get(req.Token)
	if rdbGet.Err() != nil {
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userId, _ := rdbGet.Int()
	if userId == 0 {
		// maybe fallback to postgres for token verification
		return requests.NewErrorResponse(errors.InvalidAccessTokenError())
	}
	userSlice := uniqueValues(req.InvitedUsers)
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	defer tx.Rollback()
	var foundUsers int
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(id)
		FROM users
		WHERE id = ANY($1::int[]) AND $2 != ALL(chats)
		`,
		pq.Array(userSlice),
		req.ChatId,
	).Scan(&foundUsers)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if foundUsers != len(userSlice) {
		return requests.NewErrorResponse(errors.UnableToInviteError())
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE users
		SET chats = array_append(chats, $1)
		WHERE id = ANY($2::int[])
		`,
		req.ChatId,
		pq.Array(userSlice),
	)

	if err != nil {
		fmt.Println("here")
		fmt.Println(req.ChatId, userSlice)
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	if err = tx.Commit(); err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	return requests.NewEmptyResponse(http.StatusOK)
}
