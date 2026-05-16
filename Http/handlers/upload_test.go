package handlers

import (
	"bytes"
	"fileupload/config"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidarTipo(t *testing.T) {
	cases := []struct {
		contentType string
		expect      bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"application/pdf", true},
		{"application/exe", false},
		{"text/plain", false},
	}

	for _, tc := range cases {
		result := validarTipo(tc.contentType)

		if result != tc.expect {
			t.Errorf("validarTipo(%s): expected %v, got %v",
				tc.contentType, tc.expect, result)
		}
	}
}

func TestValidUpload(t *testing.T) {
	//table driven tests
	cases := []struct {
		name           string
		nameFile       string
		content        []byte
		expectedStatus int
	}{
		{
			name:           "jpeg valido",
			nameFile:       "teste.jpg",
			content:        []byte("\xff\xd8\xff" + "conteudo fake"),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "exe invalido",
			nameFile:       "teste.exe",
			content:        []byte("MZ content fake"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "metodo fake",
			nameFile:       "teste.jpg",
			content:        []byte("\xff\xd8\xff" + "conteudo fake"),
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			var req *http.Request

			if tc.expectedStatus == http.StatusMethodNotAllowed {
				req = httptest.NewRequest(http.MethodGet, "/upload", nil)
			} else {
				var body bytes.Buffer
				writer := multipart.NewWriter(&body)

				part, _ := writer.CreateFormFile("arquivo", tc.nameFile)
				part.Write(tc.content)
				writer.Close()

				req = httptest.NewRequest(http.MethodPost, "/upload", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}

			w := httptest.NewRecorder()
			h := &UploadHandler{Config: &config.Config{
				UploadDir:     "./testdata",
				MaxUploadSize: 10 << 20,
			}}

			h.ServeUpload(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf(" %s: expected %d, got %d", tc.name, tc.expectedStatus, w.Code)
			}
		})
	}

}
