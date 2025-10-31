package telemetry

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName        string        `envconfig:"OTEL_SERVICE_NAME" required:"true"`
	OtlpEndpoint       string        `envconfig:"OTEL_OTLP_ENDPOINT" default:""`
	SamplerRatio       float64       `envconfig:"OTEL_SAMPLER_ARG"`
	BatchTimeout       time.Duration `envconfig:"OTEL_BATCH_TIMEOUT"`
	MaxQueueSize       int           `envconfig:"OTEL_MAX_QUEUE_SIZE"`
	MaxExportBatchSize int           `envconfig:"OTEL_MAX_EXPORT_BATCH_SIZE"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load telemetry config: %w", err)
	}
	return &cfg, nil
}
