package main

import (
	"context"

	"github.com/skhanal5/payflow/internal/order/config"
	"github.com/skhanal5/payflow/internal/order/kafka"
	"github.com/skhanal5/payflow/internal/order/repository"
	"github.com/skhanal5/payflow/internal/order/server"
	"github.com/skhanal5/payflow/internal/order/service"
	"github.com/skhanal5/payflow/internal/shared/logger"
)

func main() {
	cfg := config.NewConfig()
	log := logger.InitLogger("order-service", cfg.Environment)

	db := repository.NewOrderDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)

	producer := kafka.NewOrderWriter([]string{cfg.KafkaBroker}, cfg.OrderRequestedTopic, &log)
	orderService := service.NewOrderService(db, producer, &log)

	consumer := kafka.NewOrderReader([]string{cfg.KafkaBroker}, cfg.KafkaGroupId, []string{cfg.InventoryCheckedTopic}, db, &log)
	go func() {
		log.Info().Msg("Starting inventory results consumer")
		if err := consumer.ProcessInventoryResults(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Inventory results consumer failed")
		}
	}()

	server.StartServer(cfg.GRPCPort, log, orderService)
}
