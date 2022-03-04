package config

import (
	"encoding/json"
	"log"
	"os"
)

var ConfigParams = make(map[string]interface{})

func jsonArrToStringArr(arr []interface{}) []string {
	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = v.(string)
	}
	return strArr
}

func init() {
	configFileBytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configFileBytes, &ConfigParams); err != nil {
		log.Fatal(err)
	}
	ConfigParams["worker_addrs"] = jsonArrToStringArr(ConfigParams["worker_addrs"].([]interface{}))
	ConfigParams["served_methods"] = jsonArrToStringArr(ConfigParams["served_methods"].([]interface{}))
}
