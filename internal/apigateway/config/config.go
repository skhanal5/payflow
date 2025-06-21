package config

import "github.com/skhanal5/payflow/internal/shared/env"

type Config struct {
	Environment    string
	Port           string
	OrderService   string
	ProductService string
	JWTSecretKey   string
}

func NewConfig() Config {
	return Config{
		Environment:    env.GetEnvOrPanic("ENVIRONMETN"),
		Port:           env.GetEnvOrPanic("APIGATEWAY_PORT"),
		OrderService:   env.GetEnvOrPanic("ORDER_SERVICE"),
		ProductService: env.GetEnvOrPanic("PRODUCT_SERVICE"),
		JWTSecretKey:   env.GetEnvOrPanic("JWT_SECRET_KEY"),
	}
}
