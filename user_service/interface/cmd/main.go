package main

import "net/http"

func main() {
	http.HandleFunc("/create_user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from create_user"))
	})
	http.HandleFunc("/get_token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from get_token"))
	})
	http.ListenAndServe(":8082", nil)
}
