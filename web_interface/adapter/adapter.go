package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
		httpErrBadRequest(w, fmt.Sprintf("Incorrect http method. Expected: %s, but got: %s", m, r.Method))
		return
	}

	// Reading request body to array of bytes
	reqBodyBytes, err := io.ReadAll(r.Body)
	handleError(err, w)
	r.Body.Close()

	// Checking errors
	// Validating that request body is json encoded
	if !validateBodyFormat(reqBodyBytes) {
		httpErrBadRequest(w, "Expected json encoded data")
		return
	}

	// After verifying check the required fields
	jsonMap := make(map[string]interface{})
	json.Unmarshal(reqBodyBytes, &jsonMap)
	for _, field := range config.ConfigParams["request_required_fields"].(map[string][]string)[apiMethodName] {
		found := false
		for k := range jsonMap {
			if k == field {
				found = true
				break
			}
		}
		if !found {
			httpErrBadRequest(w, fmt.Sprintf("No required key in request body: %s", field))
			return
		}
	}

	// If the access token is given, verify an access token
	if token, ok := jsonMap["token"]; ok && !validateToken(token.(string)) {
		httpErrBadRequest(w, "Invalid api token. Check the sent token or try to revoke the token")
		return
	}

	// Routing the query to needed service and returning response
	adaptedRequest, err := http.NewRequest(r.Method, service_addr+apiMethodName, strings.NewReader(string(reqBodyBytes)))
	r.Body.Close()
	handleError(err, w)
	response, err := client.Do(adaptedRequest)
	handleError(err, w)
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", response.Header.Get("Content-Length"))
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
	response.Body.Close()
}
