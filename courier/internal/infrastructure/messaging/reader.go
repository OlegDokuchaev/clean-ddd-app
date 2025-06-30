package messaging

import "github.com/segmentio/kafka-go"

func NewCourierCmdReader(config *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Address},
		GroupID: config.CourierCmdConsumerGroupID,
		Topic:   config.CourierCmdTopic,
	})
}
