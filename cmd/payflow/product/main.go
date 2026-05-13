package main

import (
	"context"

	kafkaclient "github.com/segmentio/kafka-go"

	"github.com/skhanal5/payflow/internal/product/config"
	"github.com/skhanal5/payflow/internal/product/kafka"
	"github.com/skhanal5/payflow/internal/product/repository"
	"github.com/skhanal5/payflow/internal/product/server"
	"github.com/skhanal5/payflow/internal/product/service"
	"github.com/skhanal5/payflow/internal/shared/logger"
)

func main() {
	cfg := config.NewConfig()
	log := logger.InitLogger("product-service", cfg.Environment)

	db := repository.NewProductDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)

	reader := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		GroupID: cfg.KafkaGroupId,
		Topic:   cfg.OrderRequestedTopic,
	})
	writer := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(cfg.KafkaBroker),
		Topic:    cfg.InventoryCheckedTopic,
		Balancer: &kafkaclient.LeastBytes{},
	}

	processor := kafka.NewInventoryProcessor(reader, writer, db, &log)
	go func() {
		log.Info().Msg("Starting inventory processor")
		if err := processor.Start(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Inventory processor failed")
		}
	}()

	productService := service.NewProductService(db)
	server.StartServer(cfg.GRPCPort, log, productService)
}
