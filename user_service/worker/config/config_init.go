package config

import (
	"encoding/json"
	"log"
	"os"
)

type UserServiceWorkerConfig struct {
	Port                     int      `json:"port"`
	ServedMethods            []string `json:"served_methods"`
	RedisTokenValidationAddr string   `json:"redis_token_validation_addr"`
}

var Config = new(UserServiceWorkerConfig)

func init() {
	confFileBytes, err := os.ReadFile("config.json")

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(confFileBytes, Config); err != nil {
		log.Fatal(err)
	}
}
