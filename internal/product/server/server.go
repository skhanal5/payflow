package server

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rs/zerolog"

	pb "github.com/skhanal5/payflow/gen/go/product"
	productservice "github.com/skhanal5/payflow/internal/product/service"
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

func StartServer(grpcPort string, logger zerolog.Logger, productService *productservice.ProductService) {

    lis, err := net.Listen("tcp", grpcPort)
    if err != nil {
        logger.Fatal().Err(err).Msgf("Product Server: Failed to listen on %s", grpcPort)
    }

    opts := []logging.Option{
        logging.WithLevels(customFunc),
    }

    s := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            logging.UnaryServerInterceptor(interceptorLogger(logger), opts...),
        ),
    )

    pb.RegisterProductServiceServer(s, productService)
    reflection.Register(s)

    logger.Info().Msgf("Product Service (gRPC) listening on %s", grpcPort)
    if err := s.Serve(lis); err != nil {
        logger.Fatal().Err(err).Msg("Product Server: Failed to serve gRPC")
    }
}