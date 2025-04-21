package main

import (
	"context"
	"jwt-service/internal/adapter/repositories"
	"jwt-service/internal/handler"
	"jwt-service/internal/usecases"
	"jwt-service/pkg/database"
	"jwt-service/pkg/logger"
	"strings"

	"net/http"
	"os"
	"time"
)

func main() {

	log := logger.NewLogger(os.Getenv("LOG_LEVEL"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dbpool, err := database.InitDB(ctx)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}
	defer dbpool.Close()

	sessionRepo := repositories.NewSessionRepository(dbpool)
	emailSender := usecases.MockEmailSender{}
	authUC := usecases.NewAuthUseCase(sessionRepo, os.Getenv("JWT_SECRET"), 7*24*time.Hour, &emailSender)

	authHandler := handler.NewAuthHandler(authUC)

	http.HandleFunc("/auth", authHandler.HandleAuth)
	http.HandleFunc("/refresh", authHandler.HandleRefresh)

	port := strings.TrimSpace(os.Getenv("SERVER_PORT"))
	if port == "" {
		port = "8080"
	}
	log.Println("Start server on: ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.WithError(err).Fatal("Error starting server")
	}
}
