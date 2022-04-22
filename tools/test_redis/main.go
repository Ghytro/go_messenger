package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	result, err := rdb.Get("velzn4ctnr2k6xkxygfrnk7tplgnrgecsdhnp61ecbdjsw1ves").Int()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
