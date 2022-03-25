package user_actions

import (
	"context"
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
	userDataDB, err = sql.Open("postgres", "host=localhost user=postgres password=123123 dbname=message_data sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	tx, err := userDataDB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS chat_data (
			id SERIAL PRIMARY KEY NOT NULL,
			name VARCHAR(20) NOT NULL,
			avatar_url VARCHAR(2048),
			members INTEGER [] NOT NULL DEFAULT '{}',
			admin_id INTEGER NOT NULL,
			is_public BOOL NOT NULL
		);
	`)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
