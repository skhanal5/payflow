package main

import (
	"context"

	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/product/config"
	"github.com/skhanal5/payflow/internal/product/kafka"
	"github.com/skhanal5/payflow/internal/product/repository"
	"github.com/skhanal5/payflow/internal/shared"
)

func main() {
	cfg := config.NewConfig()
	db := repository.NewProductDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	reader := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		GroupID: cfg.KafkaGroupId,
		Topic:   cfg.OrderRequestedTopic,
	})
	writer := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(cfg.KafkaBroker),
		Balancer: &kafkaclient.LeastBytes{},
	}
	logger := shared.InitLogger("product-service", cfg.Environment)
	processor := kafka.NewInventoryProcessor(
		cfg.InventoryCheckedTopic,
		reader,
		writer,
		db,
		&logger,
	)

	processor.HandleIncomingOrder(context.Background())
}
