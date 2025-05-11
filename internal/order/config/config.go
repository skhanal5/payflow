package config

type Config struct {
	KafkaBroker string
	KafkaGroupId string
	OrderTopic string
	PaymentTopic string
	InventoryTopic string
	DBHost string
	DBUser string
	DBPassword string
	DBPort string
}