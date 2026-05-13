package kafka

import (
	"context"

	"github.com/rs/zerolog"
	kafkaclient "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/skhanal5/payflow/gen/go/events"
	"github.com/skhanal5/payflow/internal/order/repository"
)

type OrderConsumer interface {
	ProcessInventoryResults(ctx context.Context) error
}

type OrderReader struct {
	reader *kafkaclient.Reader
	repo   repository.OrderRepository
	logger *zerolog.Logger
}

func NewOrderReader(brokers []string, groupID string, topics []string, repo repository.OrderRepository, logger *zerolog.Logger) *OrderReader {
	r := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
	})
	return &OrderReader{
		reader: r,
		repo:   repo,
		logger: logger,
	}
}

func (r *OrderReader) ProcessInventoryResults(ctx context.Context) error {
	for {
		message, err := r.reader.ReadMessage(ctx)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to read message")
			continue
		}

		var event events.InventoryCheckedEvent
		if err := proto.Unmarshal(message.Value, &event); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal InventoryCheckedEvent")
			continue
		}

		status := "CONFIRMED"
		if !event.AllSucceeded {
			status = "FAILED"
			for _, res := range event.Results {
				if !res.Success {
					r.logger.Warn().
						Str("order_id", event.OrderId).
						Str("product_id", res.ProductId).
						Str("reason", res.Reason).
						Msg("Item inventory check failed")
				}
			}
		}

		if err := r.repo.UpdateOrderStatus(ctx, event.OrderId, status); err != nil {
			r.logger.Error().Err(err).
				Str("order_id", event.OrderId).
				Str("status", status).
				Msg("Failed to update order status")
		}
	}
}
