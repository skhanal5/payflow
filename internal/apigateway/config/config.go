package config

import "github.com/skhanal5/payflow/internal/shared/env"

type Config struct {
	Environment    string
	Port           string
	OrderService   string
	ProductService string
	AuthService    string
	JWTSecretKey   string
}

func NewConfig() Config {
	return Config{
		Environment:    env.GetEnvOrPanic("ENVIRONMENT"),
		Port:           env.GetEnvOrPanic("APIGATEWAY_PORT"),
		OrderService:   env.GetEnvOrPanic("ORDER_SERVICE"),
		ProductService: env.GetEnvOrPanic("PRODUCT_SERVICE"),
		AuthService:    env.GetEnvOrPanic("AUTH_SERVICE"),
		JWTSecretKey:   env.GetEnvOrPanic("JWT_SECRET_KEY"),
	}
}
