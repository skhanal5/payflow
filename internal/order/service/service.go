package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/skhanal5/payflow/gen/go/order"
	"github.com/skhanal5/payflow/internal/order/repository"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	repo   repository.OrderRepository
	logger *zerolog.Logger
}

func NewOrderService(repo repository.OrderRepository, logger *zerolog.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	id, err := generateID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate order id: %v", err)
	}

	var items []repository.OrderItem
	for _, item := range req.Items {
		items = append(items, repository.OrderItem{
			OrderId:   id,
			ProductId: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	order := repository.Order{
		OrderId:         id,
		UserId:          req.UserId,
		Status:          "PENDING",
		ShippingAddress: req.ShippingAddress,
		OrderItems:      items,
	}

	created, err := s.repo.InsertOrder(ctx, &order)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to insert order")
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return toProtoOrder(created), nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.repo.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	return toProtoOrder(order), nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.repo.ListOrdersByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	var pbOrders []*pb.Order
	for i := range orders {
		pbOrders = append(pbOrders, toProtoOrder(&orders[i]))
	}
	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func toProtoOrder(o *repository.Order) *pb.Order {
	var items []*pb.OrderItem
	for _, item := range o.OrderItems {
		items = append(items, &pb.OrderItem{
			ProductId: item.ProductId,
			Price:     item.Price,
			Quantity:  int32(item.Quantity),
		})
	}

	return &pb.Order{
		Id:        o.OrderId,
		UserId:    o.UserId,
		Items:     items,
		Status:    o.Status,
		OrderDate: timestamppb.New(o.CreatedAt),
	}
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate ID: %w", err)
	}
	ts := time.Now().UnixMilli()
	return fmt.Sprintf("ORD_%d_%s", ts, hex.EncodeToString(b)), nil
}
