package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type GetFileRequest struct {
	Token  string `json:"token"`
	FileId string `json:"file_id"`
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	req := new(GetFileRequest)
	if err := json.Unmarshal(reqBytes, req); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateToken(req.Token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	file, err := os.Open("files/" + req.FileId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	io.Copy(w, file)
	defer file.Close()
}
