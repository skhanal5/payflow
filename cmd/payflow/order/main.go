// package main

// import (
// 	"log"
// 	"net"

// 	"github.com/skhanal5/payflow/internal/order/config"
// 	"github.com/skhanal5/payflow/internal/order/handler"
// 	"github.com/skhanal5/payflow/internal/order/kafka"
// 	"github.com/skhanal5/payflow/internal/order/repository"
// 	"github.com/skhanal5/payflow/internal/shared"
// 	"google.golang.org/grpc"
// )

// func main() {

// 	cfg := config.NewConfig()
// 	logger := shared.InitLogger("order-service", cfg.Environment)
// 	db := repository.NewOrderDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
// 	consumer := kafka.NewOrderReader([]string{cfg.KafkaBroker}, cfg.KafkaGroupId, []string{cfg.InventoryCheckedTopic, cfg.PaymentTopic}, &logger)
// 	producer := kafka.NewOrderWriter(cfg, &logger)
// 	orderHandler := handler.NewOrderHandler(db, consumer, producer, &logger)

// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	proto.RegisterOrderServiceServer(grpcServer, orderHandler)

// 	log.Println("gRPC server is listening on :50051")
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("Failed to serve: %v", err)
// 	}
// }
