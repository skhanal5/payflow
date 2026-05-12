package kafka

import (
	"context"

	"github.com/rs/zerolog"
	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/order/config"
	"github.com/skhanal5/payflow/internal/order/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type OrderProducer interface {
	SendOrder(ctx context.Context, order *proto.PlaceOrderRequest) error
}

type OrderWriter struct {
	writer *kafkaclient.Writer
	logger *zerolog.Logger
}

func NewOrderWriter(cfg config.Config, logger *zerolog.Logger) *OrderWriter {
	w := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(cfg.KafkaBroker),
		Topic:    cfg.OrderRequestedTopic,
		Balancer: &kafkaclient.LeastBytes{},
	}
	return &OrderWriter{
		writer: w,
		logger: logger,
	}
}

func (s *OrderWriter) SendOrder(ctx context.Context, order *proto.PlaceOrderRequest) error {
	id := order.OrderId
	value, err := protobuf.Marshal(order)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to marshal order")
		return err
	}
	message := kafkaclient.Message{
		Key:   []byte(id),
		Value: value,
	}
	return s.writer.WriteMessages(ctx, message)
}
