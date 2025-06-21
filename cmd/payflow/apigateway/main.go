package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skhanal5/payflow/internal/apigateway/config"
	"github.com/skhanal5/payflow/internal/apigateway/middleware"
	"github.com/skhanal5/payflow/internal/apigateway/router"
	"github.com/skhanal5/payflow/internal/shared/logger"
)

func main() {
	ctx := context.Background()
	logger := logger.InitLogger("apigateway", "development")
	cfg := config.NewConfig()

	httpHandler, err := router.NewRouter(ctx, logger, cfg.ProductService)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize API Gateway router")
	}

	httpHandler = middleware.LoggingMiddleware(logger)(httpHandler)
	httpHandler = middleware.AuthMiddleware(logger, cfg.JWTSecretKey)(httpHandler)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: httpHandler,
	}

	go func() {
		logger.Info().Msgf("API Gateway HTTP server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("API Gateway HTTP server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Shutting down API Gateway server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("API Gateway server forced to shutdown")
	}

	logger.Info().Msg("API Gateway server stopped.")
}