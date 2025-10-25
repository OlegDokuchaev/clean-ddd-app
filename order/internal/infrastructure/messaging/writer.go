package messaging

import (
	"github.com/segmentio/kafka-go"
)

func NewOrderCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Topic:                  config.OrderCmdTopic,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: false,
	}
}

func NewWarehouseCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Topic:                  config.WarehouseCmdTopic,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: false,
	}
}

func NewCourierCommandWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Topic:                  config.CourierCmdTopic,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: false,
	}
}

func NewOrderCommandResWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Topic:                  config.OrderCmdResTopic,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: false,
	}
}
