package kafka

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	kafkaclient "github.com/segmentio/kafka-go"
)

type OrderConsumer interface {
	ReadOrderDetails(ctx context.Context) error
}

type OrderReader struct {
	reader *kafkaclient.Reader
	logger *zerolog.Logger
}

func NewOrderReader(brokers []string, groupID string, topics []string, logger *zerolog.Logger) *OrderReader {
	r := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
	})
	return &OrderReader{
		reader: r,
		logger: logger,
	}
}

func (r *OrderReader) ReadOrderDetails(ctx context.Context) error {
	for {
		message, err := r.reader.ReadMessage(ctx)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to read message")
			return err
		}
		processMessage(message)
	}
}

func processMessage(msg kafkaclient.Message) {
	// TODO: Make this more concrete
	fmt.Print(string(msg.Value))
}
