package config

import (
	"encoding/json"
	"log"
	"os"
)

type MessageServiceInterfaceConfig struct {
	WorkerAddrs []string      `json:"worker_addrs"`
	Handlers    []HandlerData `json:"handler_data"`
}

func (c *MessageServiceInterfaceConfig) HandlerData(handlerName string) *HandlerData {
	if handlerName[0] == '/' {
		handlerName = handlerName[1:]
	}
	for _, m := range c.Handlers {
		if m.Name == handlerName {
			return &m
		}
	}
	return nil
}

type HandlerData struct {
	Name           string   `json:"name"`
	Method         string   `json:"method"`
	RequiredParams []string `json:"required_params"`
}

var Config = new(MessageServiceInterfaceConfig)

func init() {
	configFileBytes, err := os.ReadFile("../config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configFileBytes, Config); err != nil {
		log.Fatal(err)
	}
}
