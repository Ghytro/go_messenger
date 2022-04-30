package handler

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

const allowedFileIdSymbols = "0123456789qwertyuiopasdfghjklzxcvbnm"
const fileIdLength = 30

func GenerateFileId() string {
	var fileIdBuilder strings.Builder
	for i := 0; i < fileIdLength; i++ {
		fileIdBuilder.WriteByte(allowedFileIdSymbols[rand.Int31n(int32(len(allowedFileIdSymbols)))])
	}
	return fileIdBuilder.String()
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	contentType, token := r.Header.Get("Content-Type"), r.Header.Get("X-Access-Token")
	if extension == "" || token == "" || !ValidateToken(token) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileId := GenerateFileId()
	for _, err := os.Stat("files/" + fileId); !errors.Is(err, os.ErrNotExist); _, err = os.Stat("files/" + fileId) {
		fileId = GenerateFileId()
	}
	file, err := os.Create(fileId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := fileDataDB.Exec(
		`INSERT INTO
		file_storages (name, storage_addr, content_type)
		VALUES ($1, $2, $3)`,
		fileId,
		Config.InstanceAddr,
		contentType,
	); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(file, r.Body)
	w.Write([]byte(fileId))
	w.WriteHeader(http.StatusOK)
}
