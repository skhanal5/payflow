package main

import (
	"context"

	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/inventory/config"
	"github.com/skhanal5/payflow/internal/inventory/kafka"
	"github.com/skhanal5/payflow/internal/inventory/repository"
)

func main() {
	cfg := config.NewConfig()
	db := repository.NewInventoryDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPassword)
	reader := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		GroupID: cfg.KafkaGroupId,
		Topic:   cfg.OrderTopic,
	})
	writer := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(cfg.KafkaBroker),
		Balancer: &kafkaclient.LeastBytes{},
	}
	logger := config.InitLogger(cfg.Environment)
	processor := kafka.NewInventoryProcessor(
		cfg.ReservationTopic,
		cfg.FailureTopic,
		reader,
		writer,
		db,
		&logger,
	)

	processor.HandleIncomingOrder(context.Background())
}
