package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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

	dst, err := os.Create("./uploads/" + header.Filename)
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

func main() {

	os.MkdirAll("./uploads", 0755)
	dir := "./public"         //diretório para armazenamento de arquivos que será servido
	mux := http.NewServeMux() //multiplexador de rotas
	fs := http.FileServer(http.Dir(dir))

	//rotas
	mux.Handle("/", fs)
	mux.HandleFunc("/upload", HandlerUpload)

	log.Println("Servindo" + dir + "em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
