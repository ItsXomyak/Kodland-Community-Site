package handlers

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/logger"
	"auth-service/internal/services"
)

type ConfirmHandler struct {
	AuthService services.AuthService
}

func NewConfirmHandler(authService services.AuthService) *ConfirmHandler {
	return &ConfirmHandler{
		AuthService: authService,
	}
}

func (h *ConfirmHandler) ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	logger.Info("Confirm account request received")
	
	// Добавляем CORS-заголовки
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// Обрабатываем OPTIONS-запрос
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		logger.Error("Token is missing in confirmation request")
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	err := h.AuthService.ConfirmAccount(token)
	if err != nil {
		logger.Error("Account confirmation failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger.Info("Account confirmed successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account confirmed successfully",
	})
}
