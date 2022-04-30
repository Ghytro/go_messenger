package handler

import (
	"database/sql"
	"net/http"
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
	var storageAddr string
	if err := fileDataDb.QueryRow("SELECT storage_addr FROM file_storages WHERE id = $1", fileId).Scan(&storageAddr); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	http.Redirect(w, r, storageAddr+"/file/"+fileId, http.StatusTemporaryRedirect)
}
