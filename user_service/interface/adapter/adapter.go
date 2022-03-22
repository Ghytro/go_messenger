package adapter

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/loadbalancer"
	"github.com/Ghytro/go_messenger/lib/requests"
	"github.com/Ghytro/go_messenger/user_service/interface/config"
)

var loadBalancer = loadbalancer.NewLoadBalancer(config.Config.WorkerAddrs)

func handleError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SendRequest(w http.ResponseWriter, r *http.Request) {
	apiMethodName := r.URL.Path
	handlerData := config.Config.HandlerData(apiMethodName)

	// checking the http method
	if m := handlerData.Method; r.Method != m {
		requests.SendResponse(w, requests.NewErrorResponse(errors.IncorrectHttpMethodError(m, r.Method)))
		return
	}

	reqBodyBytes, err := io.ReadAll(r.Body)
	handleError(err, w)

	// verifying that we have all the necessary parameters
	jsonMap := make(map[string]interface{})
	json.Unmarshal(reqBodyBytes, &jsonMap)
	for _, param := range handlerData.RequiredParams {
		found := false
		for k := range jsonMap {
			if k == param {
				found = true
				break
			}
		}
		if !found {
			if param == "token" {
				requests.SendResponse(w, requests.NewErrorResponse(errors.NoAccessTokenError()))
			} else {
				requests.SendResponse(w, requests.NewErrorResponse(errors.MissingParameterError(param)))
			}
			return
		}
	}
	// sending request to worker
	loadBalancer.SendRequest(w, r.Method, apiMethodName, strings.NewReader(string(reqBodyBytes)))
}
