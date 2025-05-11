package handler

import (
	"context"

	"github.com/skhanal5/payflow/internal/order/kafka"
	"github.com/skhanal5/payflow/internal/order/proto"
	"github.com/skhanal5/payflow/internal/order/repository"
	"google.golang.org/grpc"
) 

type OrderServer struct {
	db repository.OrderRepository
	consumer kafka.OrderConsumer
	producer kafka.OrderProducer
}


func (s *OrderServer) PlaceOrder(ctx context.Context, in *proto.PlaceOrderRequest, opts ...grpc.CallOption) (*proto.OrderResponse, error) {
	return nil, nil
}

func (s *OrderServer) GetOrderStatus(ctx context.Context, in *proto.GetOrderStatusRequest, opts ...grpc.CallOption) (*proto.GetOrderStatusResponse, error) {
	id := in.OrderId
	order, err := s.db.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	items := []*proto.OrderItem{}
	for _, element := range order.OrderItems {
		item := proto.OrderItem{
			ProductId: element.ProductId,
			Quantity: int32(element.Quantity),
		}
		items = append(items, &item)
	}
	return &proto.GetOrderStatusResponse{
		OrderId: id,
		Status: order.Status,
		Items: items,
	}, nil
}