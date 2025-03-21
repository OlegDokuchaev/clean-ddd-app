package messaging

import "github.com/segmentio/kafka-go"

func NewWarehouseCmdReader(config *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Address},
		GroupID: config.WarehouseCmdConsumerGroupID,
		Topic:   config.WarehouseCmdTopic,
	})
}
