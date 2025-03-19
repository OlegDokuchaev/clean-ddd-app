package messaging

import (
	"github.com/segmentio/kafka-go"
)

func NewOrderCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.OrderCmdTopic,
	}
}

func NewWarehouseCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.WarehouseCmdTopic,
	}
}

func NewCourierCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.CourierCmdTopic,
	}
}

func NewOrderCommandResWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.OrderCmdResTopic,
	}
}
