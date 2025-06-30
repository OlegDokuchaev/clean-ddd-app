package messaging

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"KAFKA_ADDRESS" required:"true"`

	CourierCmdTopic           string `envconfig:"KAFKA_COURIER_COMMAND_TOPIC" required:"true"`
	CourierCmdResTopic        string `envconfig:"KAFKA_COURIER_COMMAND_RESULT_TOPIC" required:"true"`
	CourierCmdConsumerGroupID string `envconfig:"KAFKA_COURIER_COMMAND_CONSUMER_GROUP_ID" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}
	return &cfg, nil
}
