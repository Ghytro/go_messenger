package adapter

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ghytro/go_messenger/lib/errors"
	"github.com/Ghytro/go_messenger/lib/response"
	"github.com/Ghytro/go_messenger/web_interface/config"
	"github.com/go-redis/redis"
)

var client = http.Client{}
var rdb = redis.NewClient(&redis.Options{
	Addr:     config.ConfigParams["redis_token_validation_addr"].(string),
	Password: "",
	DB:       0,
})

func handleError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateToken(token string) bool {
	return rdb.Get(token).Err() != redis.Nil
}

func validateBodyFormat(bodyBytes []byte) bool {
	return json.Valid(bodyBytes)
}

func httpErrBadRequest(w http.ResponseWriter, errorMessage string) {
	http.Error(w, errorMessage, http.StatusBadRequest)
}

func RequestToService(service_addr string, w http.ResponseWriter, r *http.Request) {
	// checking the query method
	apiMethodName := r.URL.Path
	if m := config.ConfigParams["handler_data"].(map[string]config.HandlerData)[apiMethodName].Method; r.Method != m {
		response.SendResponse(w, response.NewErrorResponse(errors.IncorrectHttpMethodError(m, r.Method)))
		return
	}

	// Reading request body to array of bytes
	reqBodyBytes, err := io.ReadAll(r.Body)
	handleError(err, w)
	r.Body.Close()

	// Checking errors
	// Validating that request body is json encoded
	if !validateBodyFormat(reqBodyBytes) {
		response.SendResponse(w, response.NewErrorResponse(errors.JsonValidationError()))
		return
	}

	// After verifying check the required fields
	jsonMap := make(map[string]interface{})
	json.Unmarshal(reqBodyBytes, &jsonMap)
	for _, param := range config.ConfigParams["handler_data"].(map[string]config.HandlerData)[apiMethodName].RequiredParams {
		found := false
		for k := range jsonMap {
			if k == param {
				found = true
				break
			}
		}
		if !found {
			if param == "token" {
				response.SendResponse(w, response.NewErrorResponse(errors.NoAccessTokenError()))
			} else {
				response.SendResponse(w, response.NewErrorResponse(errors.MissingParameterError(param)))
			}
			return
		}
	}

	// If the access token is given, verify an access token
	if token, ok := jsonMap["token"]; ok && !validateToken(token.(string)) {
		response.SendResponse(w, response.NewErrorResponse(errors.InvalidAccessTokenError()))
		return
	}

	// Routing the query to needed service and returning response
	adaptedRequest, err := http.NewRequest(r.Method, service_addr+apiMethodName, strings.NewReader(string(reqBodyBytes)))
	r.Body.Close()
	if err != nil {
		handleError(err, w)
		return
	}
	response, err := client.Do(adaptedRequest)
	if err != nil {
		handleError(err, w)
		return
	}
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", response.Header.Get("Content-Length"))
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
	response.Body.Close()
}
