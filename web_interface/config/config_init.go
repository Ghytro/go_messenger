package config

import (
	"encoding/json"
	"log"
	"os"
)

type HandlerData struct {
	Name           string   `json:"name"`
	Method         string   `json:"method"`
	RequiredParams []string `json:"required_params"`
}

var ConfigParams = make(map[string]interface{})

func init() {
	confFileBytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(confFileBytes, &ConfigParams)
	if err != nil {
		log.Fatal(err)
	}

	// Converting interface{} to HandlerData
	newHandlerData := make(map[string]HandlerData)
	for k, v := range ConfigParams["handler_data"].(map[string]interface{}) {
		jsonArray := v.(map[string]interface{})["required_params"].([]interface{})
		strArray := make([]string, len(jsonArray))
		for i, val := range jsonArray {
			strArray[i] = val.(string)
		}
		newHandlerData[k] = HandlerData{
			v.(map[string]interface{})["name"].(string),
			v.(map[string]interface{})["method"].(string),
			strArray,
		}
	}
	ConfigParams["handler_data"] = newHandlerData
}
