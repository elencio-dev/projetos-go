package main

import (
	"fileupload/config"
	"fileupload/handlers"
	"fileupload/middleware"
	"log"
	"net/http"
	"os"
)

func main() {

	ctg := config.New()
	uploadHandler := &handlers.UploadHandler{Config: ctg}
	filesHandler := &handlers.FilesHandler{Config: ctg}
	os.MkdirAll(ctg.UploadDir, 0755)
	//diretório para armazenamento de arquivos que será servido
	mux := http.NewServeMux() //multiplexador de rotas
	fs := http.FileServer(http.Dir(ctg.PublicDir))

	//rotas
	mux.Handle("/", fs)
	//rotas protegidas
	mux.HandleFunc("/upload", middleware.BasicAuth(uploadHandler.ServeUpload, ctg.Usuarios))
	mux.HandleFunc("/arquivos", middleware.BasicAuth(filesHandler.HandlerUploadFiles, ctg.Usuarios))

	log.Println("Servindo" + ctg.PublicDir + "em http://localhost:8080")
	log.Fatal(http.ListenAndServe(ctg.Port, mux))
}
