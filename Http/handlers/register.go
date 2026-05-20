package handlers

import (
	"encoding/json"
	"fileupload/config"
	"fileupload/store"
	"net/http"
)

type RegisterHandler struct {
	Config    *config.Config
	UserStore *store.UserStore
}

type RegisterRequest struct {
	Usuario string `json:"usuario"`
	Email   string `json:"email"`
	Senha   string `json:"senha"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func (h *RegisterHandler) ServeRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Requisição Invalida", http.StatusBadRequest)
		return
	}

	err = h.UserStore.CriarUsuario(req.Usuario, req.Email, req.Senha)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"mensagem": "Usuario Criado com Sucesso!",
	})

}
