package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path := r.URL.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	fileId := path[strings.LastIndex(path, "/")+1:]
	var contentType string
	if err := fileDataDB.QueryRow("SELECT content_type FROM file_storages WHERE id = $1", fileId).Scan(&contentType); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file, err := os.Open("files/" + fileId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	w.Header().Add("Content-Type", contentType)
	io.Copy(w, file)
}
