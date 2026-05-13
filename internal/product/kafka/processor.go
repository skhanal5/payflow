package kafka

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	kafkaclient "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/skhanal5/payflow/gen/go/events"
	"github.com/skhanal5/payflow/internal/product/repository"
)

type InventoryProcessor struct {
	reader *kafkaclient.Reader
	writer *kafkaclient.Writer
	repo   repository.ProductRepository
	logger *zerolog.Logger
}

func NewInventoryProcessor(reader *kafkaclient.Reader, writer *kafkaclient.Writer, repo repository.ProductRepository, logger *zerolog.Logger) *InventoryProcessor {
	return &InventoryProcessor{
		reader: reader,
		writer: writer,
		repo:   repo,
		logger: logger,
	}
}

func (p *InventoryProcessor) Start(ctx context.Context) error {
	for {
		message, err := p.reader.ReadMessage(ctx)
		if err != nil {
			p.logger.Error().Err(err).Msg("Failed to read order event")
			continue
		}

		var event events.OrderPlacedEvent
		if err := proto.Unmarshal(message.Value, &event); err != nil {
			p.logger.Error().Err(err).Msg("Failed to unmarshal OrderPlacedEvent")
			continue
		}

		p.processOrder(ctx, &event)
	}
}

func (p *InventoryProcessor) processOrder(ctx context.Context, event *events.OrderPlacedEvent) {
	allSucceeded := true
	var results []*events.ItemCheckResult

	for _, item := range event.Items {
		success := true
		reason := ""

		_, err := p.repo.UpdateProduct(ctx, item.ProductId, item.Quantity)
		if err != nil {
			p.logger.Error().Err(err).
				Str("order_id", event.OrderId).
				Str("product_id", item.ProductId).
				Msg("Failed to update product")
			success = false
			reason = fmt.Sprintf("insufficient stock for product %s", item.ProductId)
		}

		if !success {
			allSucceeded = false
		}

		results = append(results, &events.ItemCheckResult{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Success:   success,
			Reason:    reason,
		})
	}

	p.emitOrderResult(ctx, event.OrderId, allSucceeded, results)
}

func (p *InventoryProcessor) emitOrderResult(ctx context.Context, orderID string, allSucceeded bool, results []*events.ItemCheckResult) {
	result := &events.InventoryCheckedEvent{
		OrderId:      orderID,
		AllSucceeded: allSucceeded,
		Results:      results,
	}

	value, err := proto.Marshal(result)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to marshal InventoryCheckedEvent")
		return
	}

	message := kafkaclient.Message{
		Key:   []byte(orderID),
		Value: value,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		p.logger.Error().Err(err).
			Str("order_id", orderID).
			Msg("Failed to emit inventory check result")
	}
}
