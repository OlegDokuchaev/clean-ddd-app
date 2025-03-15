package messaging

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"KAFKA_ADDRESS" required:"true"`

	OrderCommandTopic     string `envconfig:"KAFKA_ORDER_COMMAND_TOPIC" required:"true"`
	WarehouseCommandTopic string `envconfig:"KAFKA_WAREHOUSE_COMMAND_TOPIC" required:"true"`
	CourierCommandTopic   string `envconfig:"KAFKA_COURIER_COMMAND_TOPIC" required:"true"`

	OrderCommandConsumerGroupID string `envconfig:"KAFKA_ORDER_COMMAND_CONSUMER_GROUP_ID" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}
	return &cfg, nil
}
