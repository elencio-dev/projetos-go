package config

import "log/slog"

type Config struct {
	Port          string
	UploadDir     string
	PublicDir     string
	MaxUploadSize int64
	Usuarios      map[string]string
	Logger        *slog.Logger
}

func New() *Config {
	return &Config{
		Port:          ":8080",
		UploadDir:     "./uploads",
		PublicDir:     "./public",
		MaxUploadSize: 10 << 20, //10MB
		Usuarios: map[string]string{
			"admin": "123456",
			"user":  "senha",
		},
		Logger: NewLogger(),
	}
}
