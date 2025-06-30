package messaging

import "github.com/segmentio/kafka-go"

func NewCourierCmdResWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(config.Address),
		Topic: config.CourierCmdResTopic,
	}
}
