package config

import (
	"encoding/json"
	"log"
	"os"
)

type WebInterfaceConfig struct {
	ChatServiceAddr          string `json:"chat_service_addr"`
	UserServiceAddr          string `json:"user_service_addr"`
	FileStorageServiceAddr   string `json:"file_storage_service_addr"`
	RedisTokenValidationAddr string `json:"redis_token_validation_addr"`
}

var Config = new(WebInterfaceConfig)

func init() {
	confFileBytes, err := os.ReadFile("../config/config.json")

	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(confFileBytes, Config); err != nil {
		log.Fatal(err)
	}
}
