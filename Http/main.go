package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func validarTipo(contentType string) bool {
	tiposPermitidos := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}

	permitido, existe := tiposPermitidos[contentType]
	return existe && permitido
}

func HandlerUpload(w http.ResponseWriter, r *http.Request) {

	//primeiro verificamos o metodo usado
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) //limitar o tamnho do upload do arquivo
	err := r.ParseMultipartForm(10 << 20)

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

	name := filepath.Base(header.Filename)

	dst, err := os.Create("./uploads/" + name)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo: ", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Erro ao copiar os dados: ", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Arquivo %s enviado com sucesso", header.Filename)

}

func basicAuth(handler http.HandlerFunc, usuarios map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restrito"`)
			http.Error(w, "Autenticação necessaria", http.StatusUnauthorized)
			return
		}

		if senha, existe := usuarios[user]; !existe || senha != pass {
			http.Error(w, "Credencias Inválidas", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func main() {

	os.MkdirAll("./uploads", 0755)
	dir := "./public"         //diretório para armazenamento de arquivos que será servido
	mux := http.NewServeMux() //multiplexador de rotas
	fs := http.FileServer(http.Dir(dir))

	usuarios := map[string]string{ //criando usuarios
		"admin": "123456",
		"user":  "senha",
	}

	//rotas
	mux.Handle("/", fs)
	mux.HandleFunc("/upload", basicAuth(HandlerUpload, usuarios))

	log.Println("Servindo" + dir + "em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
