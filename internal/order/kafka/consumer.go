package kafka

import (
	"context"
	"fmt"
	kafkaclient "github.com/segmentio/kafka-go"
)

type OrderConsumer interface {
	ReadOrderDetails(ctx context.Context) error
}

type OrderReader struct {
	reader *kafkaclient.Reader
}

func NewOrderReader(brokers []string, groupID string, topics []string) *OrderReader {
	r := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
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
