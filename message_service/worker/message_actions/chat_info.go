package message_actions

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/go-redis/redis"
	"github.com/lib/pq"
)

func ChatInfo(chatInfoRequest requests.Request) requests.Response {
	req := chatInfoRequest.(*requests.ChatInfoRequest)
	userId, err := redisClient.Get(req.Token).Int()
	if err != nil {
		if err == redis.Nil {
			return requests.NewErrorResponse(errors.InvalidAccessTokenError())
		}
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	rows, err := messageDataDB.Query(
		`SELECT id, name, avatar_url, is_public,
		CASE WHEN is_public OR $2 = ANY(members) THEN members ELSE NULL END,
		CASE WHEN is_public OR $2 = ANY(members) THEN admin_id ELSE NULL END
		FROM chat_data
		WHERE id = ANY($1::int[])`,
		pq.Array(req.ChatIds),
		userId,
	)
	if err != nil {
		log.Println(err)
		return requests.NewEmptyResponse(http.StatusInternalServerError)
	}
	response := requests.ChatInfoResponse{make([]requests.ChatInfo, 0, len(req.ChatIds))}
	for rows.Next() {
		var chatInfo requests.ChatInfo
		if err := rows.Scan(
			&chatInfo.Id,
			&chatInfo.Name,
			&chatInfo.AvatarUrl,
			&chatInfo.IsPublic,
			&chatInfo.Members,
			&chatInfo.AdminId,
		); err != nil {
			log.Println(err)
			return requests.NewEmptyResponse(http.StatusInternalServerError)
		}
		response.Chats = append(response.Chats, chatInfo)
	}
	return response
}
