package messaging

import (
	"github.com/segmentio/kafka-go"
)

func NewOrderCommandReader(config *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Address},
		GroupID: config.OrderCmdConsumerGroupID,
		Topic:   config.OrderCmdTopic,
	})
}

func NewWarehouseCommandResultReader(config *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Address},
		GroupID: config.WarehouseCmdResConsumerGroupID,
		Topic:   config.WarehouseCmdResTopic,
	})
}

func NewCourierCommandResultReader(config *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Address},
		GroupID: config.CourierCmdResConsumerGroupID,
		Topic:   config.CourierCmdResTopic,
	})
}
