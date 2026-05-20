package main

import (
	"context"
	"fileupload/config"
	"fileupload/db"
	"fileupload/handlers"
	"fileupload/middleware"
	"fileupload/store"
	"log"
	"net/http"
	"os"
)

func main() {
	ctg := config.New()

	//Connect DB

	client, err := db.Connect(ctg.MongoURI)
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}

	defer client.Disconnect(context.Background())

	database := client.Database(ctg.MongoDB)
	userStore := store.NovoUsuario(database)

	uploadHandler := &handlers.UploadHandler{Config: ctg}
	filesHandler := &handlers.FilesHandler{Config: ctg}
	loginHandler := handlers.LoginHandler{Config: ctg,
		UserStore: userStore}
	registerHandler := &handlers.RegisterHandler{
		Config:    ctg,
		UserStore: userStore,
	}
	os.MkdirAll(ctg.UploadDir, 0755)

	rl := middleware.NewRateLimiter(5)

	//diretório para armazenamento de arquivos que será servido
	mux := http.NewServeMux() //multiplexador de rotas
	fs := http.FileServer(http.Dir(ctg.PublicDir))

	//rotas
	mux.Handle("/", fs)
	//rotas protegidas
	mux.HandleFunc("/upload", middleware.JwtAuth(
		middleware.RateLimit(rl, uploadHandler.ServeUpload),
	))
	mux.HandleFunc("/arquivos", middleware.JwtAuth(filesHandler.HandlerUploadFiles))
	mux.HandleFunc("/login", loginHandler.ServeLogin)
	mux.HandleFunc("/register", registerHandler.ServeRegister)

	log.Println("Servindo" + ctg.PublicDir + " em http://localhost:8080")
	log.Fatal(http.ListenAndServe(ctg.Port, mux))
}
