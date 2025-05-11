package config

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
}

func NewConfig() Config {
	return Config{
		KafkaBroker: GetEnvOrPanic("KAFKA_BROKER"),
		KafkaGroupId: GetEnvOrPanic("KAFKA_GROUPID"),
		OrderTopic: GetEnvOrPanic("ORDER_TOPIC"),
		PaymentTopic: GetEnvOrPanic("PAYMENT_TOPIC"),
		InventoryTopic: GetEnvOrPanic("INVENTORY_TOPIC"),
		DBHost: GetEnvOrPanic("DATABASE_HOST"),
		DBUser: GetEnvOrPanic("DATABASE_USER"),
		DBPassword: GetEnvOrPanic("DATABASE_PASSWORD"),
		DBPort: GetEnvOrPanic("DATABASE_PORT"),
	}
}