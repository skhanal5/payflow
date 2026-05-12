package config

import "github.com/skhanal5/payflow/internal/shared/env"

type Config struct {
	KafkaBroker           string
	KafkaGroupId          string
	OrderRequestedTopic   string
	InventoryCheckedTopic string
	DBHost                string
	DBUser                string
	DBPassword            string
	DBPort                string
	GRPCPort              string
	Environment           string
}

func NewConfig() Config {
	return Config{
		KafkaBroker:           env.GetEnvOrPanic("KAFKA_BROKER"),
		KafkaGroupId:          env.GetEnvOrPanic("KAFKA_GROUPID"),
		OrderRequestedTopic:   env.GetEnvOrPanic("ORDER_REQUESTED_TOPIC"),
		InventoryCheckedTopic: env.GetEnvOrPanic("INVENTORY_CHECKED_TOPIC"),
		DBHost:                env.GetEnvOrPanic("INVENTORY_DB_HOST"),
		DBUser:                env.GetEnvOrPanic("INVENTORY_DB_USERNAME"),
		DBPassword:            env.GetEnvOrPanic("INVENTORY_DB_PASSWORD"),
		DBPort:                env.GetEnvOrPanic("INVENTORY_DB_PORT"),
		GRPCPort:              env.GetEnvOrPanic("PRODUCT_GRPC_PORT"),
		Environment:           env.GetEnvOrPanic("ENVIRONMENT"),
	}
}
