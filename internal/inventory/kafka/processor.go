package kafka

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/skhanal5/payflow/internal/inventory/proto"
	"github.com/skhanal5/payflow/internal/inventory/repository"
	order "github.com/skhanal5/payflow/internal/order/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type InventoryManager interface {
	HandleIncomingOrder(ctx context.Context)
	ReadOrderDetails(ctx context.Context) (*order.PlaceOrderRequest, error)
	ReserveInventory(ctx context.Context, order *order.PlaceOrderRequest) error
	EmitFailure(ctx context.Context, order *order.PlaceOrderRequest) error
}

type InventoryProcessor struct {
	reservedTopic string
	failureTopic string
	reader *kafkaclient.Reader
	writer *kafkaclient.Writer
	repo repository.InventoryRepository
	logger *zerolog.Logger
}

func NewInventoryProcessor(
	reservedTopic string,
	failureTopic string,
	reader *kafkaclient.Reader,
	writer *kafkaclient.Writer,
	repo repository.InventoryRepository,
	logger *zerolog.Logger,
) *InventoryProcessor {
	return &InventoryProcessor{
		reservedTopic: reservedTopic,
		failureTopic:  failureTopic,
		reader:        reader,
		writer:        writer,
		repo:          repo,
		logger:        logger,
	}
}


func (r *InventoryProcessor) HandleIncomingOrder(ctx context.Context) {
	for {
		incomingOrder, err := r.ReadOrderDetails(ctx)
		orderID := incomingOrder.OrderId
		if err != nil {
			r.logger.Error().Err(err).Msg("Error reading order details")
			continue
		}
		for _, item := range incomingOrder.Items {
			updatedProduct, err := r.repo.UpdateProduct(ctx, item.ProductId, item.Quantity)
			r.logger.Info().Interface("product", updatedProduct).Msg("Product updated")
			if err != nil {
				r.logger.Error().Err(err).Msg("Couldn't update product due to inventory or system err")
				r.EmitFailure(ctx, orderID, item.ProductId, item.Quantity)
			}	
			r.ReserveInventory(ctx, orderID, item.ProductId, item.Quantity)
		}
	}
}

func (r *InventoryProcessor) ReadOrderDetails(ctx context.Context) (*order.PlaceOrderRequest, error) {
	message, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	var order order.PlaceOrderRequest
	err = protobuf.Unmarshal(message.Value, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}


func (s *InventoryProcessor) ReserveInventory(ctx context.Context, orderID string, productID string, quantity int32) error {
	reservation := &proto.InventoryReserved{
		OrderId:   orderID,
		ProductId: productID,
		Quantity:  quantity,
		Status:    proto.InventoryStatus_RESERVED,
		Timestamp: time.Now().UnixMilli(),
	}

	value, err := protobuf.Marshal(reservation)
	if err != nil {
		return err
	}
	message := kafkaclient.Message{
		Topic: s.reservedTopic,
		Key:   []byte(orderID),
		Value: value,
	}
	return s.writer.WriteMessages(ctx, message)
}

func (s *InventoryProcessor) EmitFailure(ctx context.Context, orderID string, productID string, quantity int32) error {
	reservation := &proto.InventoryFailed{
		OrderId:   orderID,
		ProductId: productID,
		Quantity:  quantity,
		Status:	proto.InventoryStatus_FAILED,
		Reason: proto.FailureReason_OUT_OF_STOCK,
		Timestamp: time.Now().UnixMilli(),
	}
	value, err := protobuf.Marshal(reservation)
	if err != nil {
		return err
	}
	message := kafkaclient.Message{
		Topic: s.failureTopic,
		Key:   []byte(orderID),
		Value: value,
	}
	return s.writer.WriteMessages(ctx, message)
}
