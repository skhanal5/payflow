package config

import "github.com/skhanal5/payflow/internal/utility"

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
		KafkaBroker:           utility.GetEnvOrPanic("KAFKA_BROKER"),
		KafkaGroupId:          utility.GetEnvOrPanic("KAFKA_GROUPID"),
		OrderRequestedTopic:   utility.GetEnvOrPanic("ORDER_REQUESTED_TOPIC"),
		InventoryCheckedTopic: utility.GetEnvOrPanic("INVENTORY_CHECKED_TOPIC"),
		DBHost:                utility.GetEnvOrPanic("INVENTORY_DB_HOST"),
		DBUser:                utility.GetEnvOrPanic("INVENTORY_DB_USERNAME"),
		DBPassword:            utility.GetEnvOrPanic("INVENTORY_DB_PASSWORD"),
		DBPort:                utility.GetEnvOrPanic("INVENTORY_DB_PORT"),
		Environment:           utility.GetEnvOrPanic("ENVIRONMENT"),
	}
}
