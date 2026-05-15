package main

import (
	"log"
	"net/http"
)

func main() {

	dir := "./public" //diretório para armazenamento de arquivos que será servido
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Println("Servindo" + dir + "em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
