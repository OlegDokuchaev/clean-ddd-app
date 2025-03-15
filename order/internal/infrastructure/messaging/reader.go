package messaging

import (
	"github.com/segmentio/kafka-go"
)

func NewOrderCommandReader(config *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Address},
		GroupID: config.OrderCommandConsumerGroupID,
		Topic:   config.OrderCommandTopic,
	})
}
