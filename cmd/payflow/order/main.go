package main

import (
	"github.com/skhanal5/payflow/internal/order/config"
	"github.com/skhanal5/payflow/internal/order/repository"
	"github.com/skhanal5/payflow/internal/order/server"
	"github.com/skhanal5/payflow/internal/order/service"
	"github.com/skhanal5/payflow/internal/shared/logger"
)

func main() {
	cfg := config.NewConfig()
	log := logger.InitLogger("order-service", cfg.Environment)
	db := repository.NewOrderDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	orderService := service.NewOrderService(db, &log)

	server.StartServer(cfg.GRPCPort, log, orderService)
}
