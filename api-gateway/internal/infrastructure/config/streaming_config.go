package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type StreamingConfig struct {
	FileChunkSizeBytes int `envconfig:"FILE_STREAM_CHUNK_SIZE_BYTES" required:"true"`
}

func NewStreamingConfig() (*StreamingConfig, error) {
	var cfg StreamingConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load streaming config: %w", err)
	}
	return &cfg, nil
}
