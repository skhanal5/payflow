package config

import "github.com/skhanal5/payflow/internal/shared"

type Config struct {
	KafkaBroker           string
	KafkaGroupId          string
	OrderRequestedTopic   string
	InventoryCheckedTopic string
	FailureTopic          string
	DBHost                string
	DBUser                string
	DBPassword            string
	DBPort                string
	Environment           string
}

func NewConfig() Config {
	return Config{
		KafkaBroker:           shared.GetEnvOrPanic("KAFKA_BROKER"),
		KafkaGroupId:          shared.GetEnvOrPanic("KAFKA_GROUPID"),
		OrderRequestedTopic:   shared.GetEnvOrPanic("ORDER_REQUESTED_TOPIC"),
		InventoryCheckedTopic: shared.GetEnvOrPanic("INVENTORY_CHECKED_TOPIC"),
		DBHost:                shared.GetEnvOrPanic("INVENTORY_DB_HOST"),
		DBUser:                shared.GetEnvOrPanic("INVENTORY_DB_USERNAME"),
		DBPassword:            shared.GetEnvOrPanic("INVENTORY_DB_PASSWORD"),
		DBPort:                shared.GetEnvOrPanic("INVENTORY_DB_PORT"),
		Environment:           shared.GetEnvOrPanic("ENVIRONMENT"),
	}
}
