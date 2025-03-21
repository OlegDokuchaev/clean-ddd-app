package messaging

import "github.com/segmentio/kafka-go"

func NewWarehouseCmdResWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.WarehouseCmdResTopic,
	}
}

func NewProductEventWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.ProductEventTopic,
	}
}
