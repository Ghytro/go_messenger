package config

import (
	"encoding/json"
	"log"
	"os"
)

type FileStorageServiceWorkerConfig struct {
	RedisTokenValidationAddr string `json:"redis_token_validation_addr"`
}

var Config = new(FileStorageServiceWorkerConfig)

func init() {
	configBytes, err := os.ReadFile("../config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configBytes, Config); err != nil {
		log.Fatal(err)
	}
}
