package config

import "github.com/skhanal5/payflow/internal/utility"

type Config struct {
	KafkaBroker    string
	KafkaGroupId   string
	OrderTopic     string
	PaymentTopic   string
	InventoryTopic string
	DBHost         string
	DBUser         string
	DBPassword     string
	DBPort         string
	Environment    string
}

func NewConfig() Config {
	return Config{
		KafkaBroker:    utility.GetEnvOrPanic("KAFKA_BROKER"),
		KafkaGroupId:   utility.GetEnvOrPanic("KAFKA_GROUPID"),
		OrderTopic:     utility.GetEnvOrPanic("ORDER_TOPIC"),
		PaymentTopic:   utility.GetEnvOrPanic("PAYMENT_TOPIC"),
		InventoryTopic: utility.GetEnvOrPanic("INVENTORY_TOPIC"),
		DBHost:         utility.GetEnvOrPanic("DATABASE_HOST"),
		DBUser:         utility.GetEnvOrPanic("DATABASE_USER"),
		DBPassword:     utility.GetEnvOrPanic("DATABASE_PASSWORD"),
		DBPort:         utility.GetEnvOrPanic("DATABASE_PORT"),
		Environment:    utility.GetEnvOrPanic("ENVIRONMENT"),
	}
}
