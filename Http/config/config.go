package config

import "log/slog"

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
}

func New() *Config {
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
		Logger: NewLogger(),
	}
}
