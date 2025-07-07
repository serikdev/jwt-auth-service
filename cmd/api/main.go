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

	_ "jwt-service/docs"

	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           JWT Auth Service API
// @version         1.0
// @description     Выдаёт и обновляет Access/Refresh JWT.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     "Bearer <token>"

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dbpool, err := database.InitDB(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}
	defer dbpool.Close()
	log := logger.NewLogger()

	sessionRepo := repositories.NewSessionRepository(dbpool)
	emailSender := usecases.MockEmailSender{}
	authUC := usecases.NewAuthUseCase(sessionRepo, os.Getenv("JWT_SECRET"), 7*24*time.Hour, &emailSender)

	mux := http.NewServeMux()
	authHandler := handler.NewAuthHandler(authUC)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/auth", authHandler.HandleAuth)
	mux.HandleFunc("/refresh", authHandler.HandleRefresh)

	port := strings.TrimSpace(os.Getenv("SERVER_PORT"))
	if port == "" {
		port = "8080"
	}
	log.Println("Start server on:", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.WithError(err).Fatal("Error starting server")
	}
}
