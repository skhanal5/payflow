package main

import (
	"log"
	"net"

	"github.com/skhanal5/payflow/internal/order/config"
	"github.com/skhanal5/payflow/internal/order/handler"
	"github.com/skhanal5/payflow/internal/order/kafka"
	"github.com/skhanal5/payflow/internal/order/proto"
	"github.com/skhanal5/payflow/internal/order/repository"
	"google.golang.org/grpc"
)


func main() {

	cfg := config.NewConfig()
	db := repository.NewOrderDB(cfg)
	consumer := kafka.NewOrderReader(cfg)
	producer := kafka.NewOrderWriter(cfg)
	orderHandler := handler.NewOrderHandler(db, consumer, producer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterOrderServiceServer(grpcServer, orderHandler)

	log.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}