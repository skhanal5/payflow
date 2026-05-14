package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authgw "github.com/skhanal5/payflow/gen/go/auth"
	ordergw "github.com/skhanal5/payflow/gen/go/order"
	productgw "github.com/skhanal5/payflow/gen/go/product"

	"github.com/skhanal5/payflow/internal/apigateway/handler"
)

func NewRouter(ctx context.Context, logger zerolog.Logger, productAddr, orderAddr, authAddr string) (http.Handler, error) {
	gwmux := runtime.NewServeMux()

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := productgw.RegisterProductServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		productAddr,
		dialOpts,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to register product service gateway handler from endpoint %s", productAddr)
		return nil, fmt.Errorf("failed to register product service gateway: %w", err)
	}
	logger.Info().Msg("Registered product service gateway handler")

	err = ordergw.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		orderAddr,
		dialOpts,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to register order service gateway handler from endpoint %s", orderAddr)
		return nil, fmt.Errorf("failed to register order service gateway: %w", err)
	}
	logger.Info().Msg("Registered order service gateway handler")

	err = authgw.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		authAddr,
		dialOpts,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to register auth service gateway handler from endpoint %s", authAddr)
		return nil, fmt.Errorf("failed to register auth service gateway: %w", err)
	}
	logger.Info().Msg("Registered auth service gateway handler")

	mainMux := http.NewServeMux()

	mainMux.HandleFunc("/health", handler.GetHealth)

	mainMux.Handle("/", gwmux)

	return mainMux, nil
}
