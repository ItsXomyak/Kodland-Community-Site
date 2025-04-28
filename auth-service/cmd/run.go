package cmd

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"auth-service/config"
	"auth-service/internal/api/handlers"
	"auth-service/internal/api/handlers/routes"
	"auth-service/internal/logger"
	"auth-service/internal/mailer"
	"auth-service/internal/repository"
	"auth-service/internal/services"
)

func Run() {
	logger.Init()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Error("Error loading config: ", err)
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("Error connecting to database: ", err)
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewConfirmationTokenRepository(db)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db)

	smtpMailer := mailer.NewSMTPMailer(cfg)

	authService := services.NewAuthService(userRepo, tokenRepo, passwordResetTokenRepo, cfg, smtpMailer)

	authHandler := handlers.NewAuthHandler(authService)
	confirmHandler := handlers.NewConfirmHandler(authService)
	passwordResetHandler := handlers.NewPasswordResetHandler(authService)
	profileHandler := handlers.NewProfileHandler(db)

	router := routes.RegisterRoutes(authHandler, confirmHandler, passwordResetHandler, profileHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Origin", "X-Requested-With", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 часа
		Debug:            true,  // Включаем отладку для разработки
	}).Handler(router)

	logger.Info("Server starting on port: ", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, corsHandler); err != nil {
		logger.Error("Server failed: ", err)
		log.Fatalf("Server failed: %v", err)
	}
}
