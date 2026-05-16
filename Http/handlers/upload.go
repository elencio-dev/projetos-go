package handlers

import (
	"fileupload/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadHandler struct {
	Config *config.Config
}

func (h *UploadHandler) ServeUpload(w http.ResponseWriter, r *http.Request) {
	//primeiro verificamos o metodo usado
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.Config.MaxUploadSize) //limitar o tamnho do upload do arquivo
	err := r.ParseMultipartForm(h.Config.MaxUploadSize)

	if err != nil {
		http.Error(w, "Erro ao processar arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("arquivo")
	if err != nil {
		http.Error(w, "Erro ao processar arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// le os primeiros bytes para verificar o tipo real do arquivo
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Erro ao ler arquivo: ", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(buffer)

	if !validarTipo(contentType) {
		http.Error(w, "Tipo de arquivo não Permitido"+contentType, http.StatusBadRequest)
		return
	}

	file.Seek(0, io.SeekStart)

	name := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))

	dst, err := os.Create(filepath.Join(h.Config.UploadDir, name))
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo: ", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Erro ao copiar os dados: ", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Arquivo %s enviado com sucesso", name)

}

func validarTipo(contentType string) bool {
	tiposPermitidos := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}

	permitido, existe := tiposPermitidos[contentType]
	return existe && permitido
}
