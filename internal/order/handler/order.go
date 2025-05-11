package handler

import (
	"context"

	"github.com/skhanal5/payflow/internal/order/kafka"
	"github.com/skhanal5/payflow/internal/order/proto"
	"github.com/skhanal5/payflow/internal/order/repository"
)

type OrderHandler struct {
	proto.UnimplementedOrderServiceServer
	db       repository.OrderRepository
	consumer kafka.OrderConsumer
	producer kafka.OrderProducer
}

func NewOrderHandler(db repository.OrderRepository, consumer kafka.OrderConsumer, producer kafka.OrderProducer) *OrderHandler {
	return &OrderHandler{
		db:       db,
		consumer: consumer,
		producer: producer,
	}
}

func (h *OrderHandler) PlaceOrder(ctx context.Context, in *proto.PlaceOrderRequest) (*proto.OrderResponse, error) {
	order := convertToDBItem(in)
	res, err := h.db.InsertOrder(ctx, &order)
	if err != nil {
		return nil, err
	}
	return &proto.OrderResponse{
		OrderId: res.OrderId,
		Status:  res.Status,
	}, nil
}

func (h *OrderHandler) GetOrderStatus(ctx context.Context, in *proto.GetOrderStatusRequest) (*proto.GetOrderStatusResponse, error) {
	id := in.OrderId
	order, err := h.db.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	items := getResponseItems(order.OrderItems)
	return &proto.GetOrderStatusResponse{
		OrderId: id,
		Status:  order.Status,
		Items:   items,
	}, nil
}

func convertToDBItem(in *proto.PlaceOrderRequest) repository.Order {
		items := []repository.OrderItem{}
	for _, element := range in.Items {
		item := repository.OrderItem{
			OrderId:   in.OrderId,
			ProductId: element.ProductId,
			Quantity:  int(element.Quantity),
		}
		items = append(items, item)
	}

	order := repository.Order{
		OrderId:    in.OrderId,
		Status:     "Placed",
		OrderItems: items,
	}
	return order
}

func getResponseItems(orderItems []repository.OrderItem) []*proto.OrderItem {
	items := []*proto.OrderItem{} // pointer might be a bug
	for _, element := range orderItems {
		item := proto.OrderItem{
			ProductId: element.ProductId,
			Quantity:  int32(element.Quantity),
		}
		items = append(items, &item)
	}
	return items
}