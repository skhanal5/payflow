package server

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/skhanal5/payflow/gen/go/order"
	orderservice "github.com/skhanal5/payflow/internal/order/service"
)

var customFunc = logging.DefaultServerCodeToLevel

func interceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func StartServer(grpcPort string, logger zerolog.Logger, orderService *orderservice.OrderService) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Order Server: Failed to listen on %s", grpcPort)
	}

	opts := []logging.Option{
		logging.WithLevels(customFunc),
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger), opts...),
		),
	)

	pb.RegisterOrderServiceServer(s, orderService)
	reflection.Register(s)

	logger.Info().Msgf("Order Service (gRPC) listening on %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("Order Server: Failed to serve gRPC")
	}
}
