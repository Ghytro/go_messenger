// This tool is supposed to set up redis with tokens
// if the redis was shut down at some reason and backups were not saved
// the data is taken from postgres database from user service
// redis db is filled with this data

package main

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type dbRow struct {
	id    int
	token string
}

func main() {
	conn, err := sql.Open("postgres", "host=localhost user=postgres password=123123 dbname=user_data sslmode=disable")
	handleErr(err)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	rows, err := conn.Query("SELECT id, access_token FROM users")
	handleErr(err)
	pairs := make([]interface{}, 0)
	for rows.Next() {
		r := dbRow{}
		rows.Scan(&r.id, &r.token)
		pairs = append(pairs, r.token, r.id)
	}
	cmd := rdb.MSet(pairs...)
	handleErr(cmd.Err())
	defer conn.Close()
}
