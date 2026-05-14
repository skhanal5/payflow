package main

import (
	"github.com/skhanal5/payflow/internal/auth/config"
	"github.com/skhanal5/payflow/internal/auth/repository"
	"github.com/skhanal5/payflow/internal/auth/server"
	"github.com/skhanal5/payflow/internal/auth/service"
	"github.com/skhanal5/payflow/internal/shared/logger"
)

func main() {
	cfg := config.NewConfig()
	log := logger.InitLogger("auth-service", cfg.Environment)

	db := repository.NewUserDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	authService := service.NewAuthService(db, cfg.JWTSecret, &log)

	server.StartServer(cfg.GRPCPort, log, authService)
}
