package kafka

import (
	"context"

	"github.com/rs/zerolog"
	kafkaclient "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/skhanal5/payflow/gen/go/events"
)

type OrderProducer interface {
	SendOrderPlaced(ctx context.Context, event *events.OrderPlacedEvent) error
}

type OrderWriter struct {
	writer *kafkaclient.Writer
	logger *zerolog.Logger
}

func NewOrderWriter(brokers []string, topic string, logger *zerolog.Logger) *OrderWriter {
	w := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafkaclient.LeastBytes{},
	}
	return &OrderWriter{
		writer: w,
		logger: logger,
	}
}

func (w *OrderWriter) SendOrderPlaced(ctx context.Context, event *events.OrderPlacedEvent) error {
	value, err := proto.Marshal(event)
	if err != nil {
		w.logger.Error().Err(err).Msg("Failed to marshal OrderPlacedEvent")
		return err
	}
	message := kafkaclient.Message{
		Key:   []byte(event.OrderId),
		Value: value,
	}
	return w.writer.WriteMessages(ctx, message)
}
