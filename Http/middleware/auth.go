package middleware

import "net/http"

func BasicAuth(handler http.HandlerFunc, usuarios map[string]string) http.HandlerFunc {
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
