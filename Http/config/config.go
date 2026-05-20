package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Usuario struct {
	Usuario string
	Email   string
	Senha   string
}

type Config struct {
	Port          string
	UploadDir     string
	PublicDir     string
	MaxUploadSize int64
	Usuarios      map[string]Usuario
	Logger        *slog.Logger
	MongoURI      string
	MongoDB       string
}

func New() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	return &Config{
		Port:          ":8080",
		UploadDir:     "./uploads",
		PublicDir:     "./public",
		MaxUploadSize: 10 << 20, //10MB
		Usuarios: map[string]Usuario{
			"admin": {
				Usuario: "Administrador",
				Email:   "admin@gmail.com",
				Senha:   "12345Type",
			},
		},
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),
		Logger:   NewLogger(),
	}
}
