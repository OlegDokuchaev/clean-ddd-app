package messaging

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"KAFKA_ADDRESS" required:"true"`

	WarehouseCmdTopic           string `envconfig:"KAFKA_WAREHOUSE_COMMAND_TOPIC" required:"true"`
	WarehouseCmdResTopic        string `envconfig:"KAFKA_WAREHOUSE_COMMAND_RESULT_TOPIC" required:"true"`
	WarehouseCmdConsumerGroupID string `envconfig:"KAFKA_WAREHOUSE_COMMAND_CONSUMER_GROUP_ID" required:"true"`

	ProductEventTopic           string `envconfig:"KAFKA_PRODUCT_EVENT_TOPIC" required:"true"`
	ProductEventConsumerGroupID string `envconfig:"KAFKA_PRODUCT_EVENT_CONSUMER_GROUP_ID" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}
	return &cfg, nil
}
