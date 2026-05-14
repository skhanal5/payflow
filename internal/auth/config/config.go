package config

import "github.com/skhanal5/payflow/internal/shared/env"

type Config struct {
	DBHost      string
	DBUser      string
	DBPassword  string
	DBPort      string
	GRPCPort    string
	JWTSecret   string
	Environment string
}

func NewConfig() Config {
	return Config{
		DBHost:      env.GetEnvOrPanic("AUTH_DB_HOST"),
		DBUser:      env.GetEnvOrPanic("AUTH_DB_USER"),
		DBPassword:  env.GetEnvOrPanic("AUTH_DB_PASSWORD"),
		DBPort:      env.GetEnvOrPanic("AUTH_DB_PORT"),
		GRPCPort:    env.GetEnvOrPanic("AUTH_GRPC_PORT"),
		JWTSecret:   env.GetEnvOrPanic("JWT_SECRET_KEY"),
		Environment: env.GetEnvOrPanic("ENVIRONMENT"),
	}
}
