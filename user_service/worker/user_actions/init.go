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
	_, err = userDataDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY NOT NULL,
			username VARCHAR(20) UNIQUE NOT NULL,
			email VARCHAR(320) UNIQUE NOT NULL, -- max lenght of email address is defined by international standart
			password_md5_hash CHAR(32) NOT NULL,
			access_token CHAR(50) UNIQUE NOT NULL,
			bio TEXT,
			avatar_url VARCHAR(2048)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
