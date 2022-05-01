package handler

import (
	"net/http"

	"github.com/Ghytro/go_messenger/file_storage_service/interface/config"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, config.Config.StoragesAddrs[counter.Inc()]+"/upload_file", http.StatusTemporaryRedirect)
}
