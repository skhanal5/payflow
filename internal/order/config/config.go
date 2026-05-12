package config

import "github.com/skhanal5/payflow/internal/shared/env"

type Config struct {
	KafkaBroker           string
	KafkaGroupId          string
	OrderRequestedTopic   string
	PaymentTopic          string
	InventoryCheckedTopic string
	DBHost                string
	DBUser                string
	DBPassword            string
	DBPort                string
	Environment           string
}

func NewConfig() Config {
	return Config{
		KafkaBroker:           env.GetEnvOrPanic("KAFKA_BROKER"),
		KafkaGroupId:          env.GetEnvOrPanic("KAFKA_GROUPID"),
		OrderRequestedTopic:   env.GetEnvOrPanic("ORDER_REQUESTED_TOPIC"),
		InventoryCheckedTopic: env.GetEnvOrPanic("INVENTORY_CHECKED_TOPIC"),
		DBHost:                env.GetEnvOrPanic("ORDER_DB_HOST"),
		DBUser:                env.GetEnvOrPanic("ORDER_DB_USER"),
		DBPassword:            env.GetEnvOrPanic("ORDER_DB_PASSWORD"),
		DBPort:                env.GetEnvOrPanic("ORDER_DB_PORT"),
		Environment:           env.GetEnvOrPanic("ENVIRONMENT"),
	}
}
