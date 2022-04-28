package handler

import (
	"github.com/Ghytro/go_messenger/notification_service/config"
	"github.com/go-redis/redis"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     config.Config.RedisTokenValidationAddr,
	Password: "",
	DB:       0,
})
