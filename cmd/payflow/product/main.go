package main

import (
	"github.com/skhanal5/payflow/internal/product/config"
	"github.com/skhanal5/payflow/internal/product/repository"
	"github.com/skhanal5/payflow/internal/product/server"
	"github.com/skhanal5/payflow/internal/product/service"
	"github.com/skhanal5/payflow/internal/shared/logger"
)

func main() {
	cfg := config.NewConfig()
	log := logger.InitLogger("product-service", cfg.Environment)
	db := repository.NewProductDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	productService := service.NewProductService(db)

	server.StartServer(cfg.GRPCPort, log, productService)
}
