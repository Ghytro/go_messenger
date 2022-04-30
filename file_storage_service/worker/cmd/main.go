package main

import (
	"log"
	"net/http"

	"github.com/Ghytro/go_messenger/file_storage_service/worker/handler"
)

func main() {
	http.HandleFunc("/get_file", handler.GetFile)
	http.HandleFunc("/upload_file", handler.UploadFile)
	log.Println(http.ListenAndServe(":8076", nil))
}
