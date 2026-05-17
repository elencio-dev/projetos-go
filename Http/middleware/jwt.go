package middleware

import (
	"fileupload/auth"
	"net/http"
	"strings"
)

func JwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token Necessario", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		usuario, err := auth.ValidarToken(tokenString)
		if err != nil {
			http.Error(w, "Token Invalido", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-Usuario", usuario)
		next(w, r)
	}
}
