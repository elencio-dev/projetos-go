package handlers

import (
	"encoding/json"
	"fileupload/auth"
	"fileupload/config"
	"net/http"
)

type LoginHandler struct {
	Config *config.Config
}

type LoginRequest struct {
	Usuario string `json:"usuario"`
	Email   string `json:"email"`
	Senha   string `json:"senha"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *LoginHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Requisição Invalida", http.StatusBadRequest)
		return
	}

	usuario, existe := h.Config.Usuarios[req.Usuario]
	if !existe || usuario.Senha != req.Senha {
		http.Error(w, "Credencias Invalidas", http.StatusUnauthorized)
		return
	}

	token, err := auth.GerarToken(req.Usuario)
	if err != nil {
		http.Error(w, "Erro ao gerar Token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})

}
