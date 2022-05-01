package handler

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var fileDataDB *sql.DB

func init() {
	var err error
	if fileDataDB, err = sql.Open("postgres", "host=localhost user=postgres password=123123 dbname=file_data sslmode=disable"); err != nil {
		log.Fatal(err)
	}
}

func ValidateToken(token string) bool {
	return redisClient.Get(token).Err() == nil
}
