package config

import (
	"encoding/json"
	"log"
	"os"
)

type UserServiceInterfaceConfig struct {
	WorkerAddrs   []string `json:"worker_addrs"`
	ServedMethods []string `json:"served_methods"`
}

var Config = new(UserServiceInterfaceConfig)

func init() {
	configFileBytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configFileBytes, Config); err != nil {
		log.Fatal(err)
	}
}
