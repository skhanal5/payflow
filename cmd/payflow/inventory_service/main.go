package main

import (
	"context"

	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/inventory/config"
	"github.com/skhanal5/payflow/internal/inventory/kafka"
	"github.com/skhanal5/payflow/internal/inventory/repository"
	"github.com/skhanal5/payflow/internal/utility"
)

func main() {
	cfg := config.NewConfig()
	db := repository.NewInventoryDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPassword)
	reader := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		GroupID: cfg.KafkaGroupId,
		Topic:   cfg.OrderRequestedTopic,
	})
	writer := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(cfg.KafkaBroker),
		Balancer: &kafkaclient.LeastBytes{},
	}
	logger := utility.InitLogger("inventory-service", cfg.Environment)
	processor := kafka.NewInventoryProcessor(
		cfg.InventoryCheckedTopic,
		reader,
		writer,
		db,
		&logger,
	)

	processor.HandleIncomingOrder(context.Background())
}
