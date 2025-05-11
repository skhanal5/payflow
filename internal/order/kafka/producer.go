package kafka

import (
	"context"

	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/order/config"
	"github.com/skhanal5/payflow/internal/order/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type OrderProducer interface {
	SendOrderFunc(ctx context.Context, order proto.PlaceOrderRequest) error
}

type OrderWriter struct {
	writer *kafkaclient.Writer
}

func NewOrderWriter(cfg config.Config) *OrderWriter {
	w := &kafkaclient.Writer{
		Addr:     kafkaclient.TCP(cfg.KafkaBroker),
		Topic:    cfg.OrderTopic,
		Balancer: &kafkaclient.LeastBytes{},
	}
	return &OrderWriter{
		writer: w,
	}
}

func (s *OrderWriter) SendOrder(ctx context.Context, order *proto.PlaceOrderRequest) error {
	id := order.OrderId
	value, err := protobuf.Marshal(order)
	if err != nil {
		return err
	}
	message := kafkaclient.Message{
		Key:   []byte(id),
		Value: value,
	}
	return s.writer.WriteMessages(ctx, message)
}
