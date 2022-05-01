package handler

import (
	"database/sql"
	"fmt"
	"net/http"
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
	var storageAddr string
	if err := fileDataDB.QueryRow("SELECT storage_addr FROM file_storages WHERE id = $1", fileId).Scan(&storageAddr); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("here2")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	http.Redirect(w, r, storageAddr+"/file?id="+fileId, http.StatusTemporaryRedirect)
}
