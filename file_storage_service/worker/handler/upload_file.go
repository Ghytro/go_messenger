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
	extension, token := r.Header.Get("X-File-Extension"), r.Header.Get("X-Access-Token")
	if extension == "" || token == "" || len(extension) > 10 || !ValidateToken(token) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	extension = strings.ToLower(extension)
	for _, c := range extension {
		if c < 'a' || c > 'z' {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	fileName := GenerateFileId() + "." + extension
	for _, err := os.Stat("files/" + fileName); !errors.Is(err, os.ErrNotExist); _, err = os.Stat("files/" + fileName) {
		fileName = GenerateFileId() + "." + extension
	}
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(file, r.Body)
	w.Write([]byte(fileName))
	w.WriteHeader(http.StatusOK)
}
