package config

import (
	"encoding/json"
	"log"
	"os"
)

type NotificationServiceConfig struct {
	RedisTokenValidationAddr string `json:"redis_token_validation_addr"`
}

var Config = new(NotificationServiceConfig)

func init() {
	contentBytes, err := os.ReadFile("../config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(contentBytes, Config); err != nil {
		log.Fatal(err)
	}
}
