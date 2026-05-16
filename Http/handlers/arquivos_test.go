package handlers

import (
	"encoding/json"
	"fileupload/config"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUploadFiles(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/arquivos", nil)

	w := httptest.NewRecorder()

	h := &FilesHandler{Config: &config.Config{UploadDir: "./testdata"}}

	h.HandlerUploadFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected application/json, got %s", contentType)
	}

	body, _ := io.ReadAll(w.Body)

	var arquivos []string
	json.Unmarshal(body, &arquivos)

	if len(arquivos) == 0 {
		t.Error("espera arquivos got 0")
	}

}
