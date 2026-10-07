package handler

import (
	"encoding/json"
	"net/http"
	"whatsapp-store/internal/domain"
)

type AdminHandler struct {
	adminUsecase domain.AdminUsecase
}

func NewAdminHandler(u domain.AdminUsecase) *AdminHandler {
	return &AdminHandler{adminUsecase: u}
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.adminUsecase.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
