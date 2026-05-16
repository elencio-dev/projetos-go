package handlers

import (
	"encoding/json"
	"fileupload/config"
	"net/http"
	"os"
)

type FilesHandler struct {
	Config *config.Config
}

func (h *FilesHandler) HandlerUploadFiles(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(h.Config.UploadDir)
	if err != nil {
		http.Error(w, "Error ao carregar arquivos", http.StatusInternalServerError)
		return
	}

	nome := []string{}

	for _, entry := range entries {
		nome = append(nome, entry.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nome)
}
