package messaging

import (
	"github.com/segmentio/kafka-go"
)

func NewOrderCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.OrderCommandTopic,
	}
}

func NewWarehouseCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.WarehouseCommandTopic,
	}
}

func NewCourierCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.CourierCommandTopic,
	}
}
