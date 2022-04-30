package config

import (
	"encoding/json"
	"log"
	"os"
)

type FileStorageServiceInterfaceConfig struct {
	StoragesAddrs []string `json:"storages_addrs"`
}

var Config = new(FileStorageServiceInterfaceConfig)

func init() {
	configBytes, err := os.ReadFile("../config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configBytes, Config); err != nil {
		log.Fatal(err)
	}
}
