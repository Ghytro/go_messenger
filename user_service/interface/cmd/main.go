package main

import (
	"net/http"

	"github.com/Ghytro/go_messenger/lib/loadbalancer"
	"github.com/Ghytro/go_messenger/user_service/interface/config"
)

var LoadBalancer = loadbalancer.NewLoadBalancer(config.ConfigParams["worker_addrs"].([]string))

func main() {
	for _, m := range config.ConfigParams["served_methods"].([]string) {
		http.HandleFunc(m, LoadBalancer.SendRequest)
	}
	http.ListenAndServe(":8082", nil)
}
