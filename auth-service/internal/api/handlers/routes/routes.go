package routes

import (
	"github.com/gorilla/mux"

	"auth-service/internal/api/handlers"
)

func RegisterRoutes(
	authHandler *handlers.AuthHandler,
	confirmHandler *handlers.ConfirmHandler,
	passwordResetHandler *handlers.PasswordResetHandler,
	profileHandler *handlers.ProfileHandler,
) *mux.Router {
	router := mux.NewRouter()

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/confirm", confirmHandler.ConfirmAccount).Methods("GET")
	auth.HandleFunc("/password-reset-request", passwordResetHandler.RequestPasswordReset).Methods("POST")
	auth.HandleFunc("/password-reset-confirm", passwordResetHandler.ResetPassword).Methods("POST")

	// Profile routes
	profile := api.PathPrefix("/profile").Subrouter()
	profileHandler.RegisterRoutes(profile)

	return router
}
