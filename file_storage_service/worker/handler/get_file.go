package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Println("here0")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileIdValues, ok := r.URL.Query()["id"]
	if !ok || len(fileIdValues) != 1 {
		fmt.Println("here")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileId := fileIdValues[0]
	if len(fileId) != 30 {
		fmt.Println("here1")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, c := range fileId {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			fmt.Println("here1337")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	var contentType string
	if err := fileDataDB.QueryRow("SELECT content_type FROM file_storages WHERE id = $1", fileId).Scan(&contentType); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file, err := os.Open("../files/" + fileId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	w.Header().Add("Content-Type", contentType)
	io.Copy(w, file)
}
