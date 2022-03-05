package user_actions

import (
	"database/sql"
	"log"

	"github.com/Ghytro/go_messenger/user_service/worker/config"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

var userDataDB *sql.DB
var redisClient = redis.NewClient(&redis.Options{
	Addr:     config.Config.RedisTokenValidationAddr,
	Password: "",
	DB:       0,
})

func init() {
	var err error
	userDataDB, err = sql.Open("postgres", "host=localhost user=postgres password=123123 dbname=user_data sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}
