package kafka

import (
	"context"
	"fmt"

	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/order/config"
)

type OrderConsumer interface {
	ReadOrderDetails(ctx context.Context) error
}

type OrderReader struct {
	reader *kafkaclient.Reader
}

func NewOrderReader(cfg config.Config) *OrderReader {
	r := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers:     []string{cfg.KafkaBroker},
		GroupID:     cfg.KafkaGroupId,
		GroupTopics: []string{cfg.InventoryTopic, cfg.PaymentTopic},
	})
	return &OrderReader{
		reader: r,
	}
}

func (r *OrderReader) ReadOrderDetails(ctx context.Context) error {
	for {
		message, err := r.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		processMessage(message)
	}
}

func processMessage(msg kafkaclient.Message) {
	// TODO: Make this more concrete
	fmt.Print(string(msg.Value))
}
