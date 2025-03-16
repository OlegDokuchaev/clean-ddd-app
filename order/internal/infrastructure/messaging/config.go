package messaging

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"KAFKA_ADDRESS" required:"true"`

	OrderCmdTopic           string `envconfig:"KAFKA_ORDER_COMMAND_TOPIC" required:"true"`
	OrderCmdConsumerGroupID string `envconfig:"KAFKA_ORDER_COMMAND_CONSUMER_GROUP_ID" required:"true"`

	WarehouseCmdTopic              string `envconfig:"KAFKA_WAREHOUSE_COMMAND_TOPIC" required:"true"`
	WarehouseCmdResTopic           string `envconfig:"KAFKA_WAREHOUSE_COMMAND_RESULT_TOPIC" required:"true"`
	WarehouseCmdResConsumerGroupID string `envconfig:"KAFKA_WAREHOUSE_COMMAND_RESULT_CONSUMER_GROUP_ID" required:"true"`

	CourierCmdTopic              string `envconfig:"KAFKA_COURIER_COMMAND_TOPIC" required:"true"`
	CourierCmdResTopic           string `envconfig:"KAFKA_COURIER_COMMAND_RESULT_TOPIC" required:"true"`
	CourierCmdResConsumerGroupID string `envconfig:"KAFKA_COURIER_COMMAND_RESULT_CONSUMER_GROUP_ID" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}
	return &cfg, nil
}
