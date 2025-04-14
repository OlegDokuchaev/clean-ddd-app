package product

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Endpoint         string `envconfig:"MINIO_ENDPOINT" required:"true"`
	BucketName       string `envconfig:"MINIO_BUCKET_NAME" required:"true"`
	DefaultImageName string `envconfig:"MINIO_DEFAULT_IMAGE_NAME" required:"true"`
	UseSSL           bool   `envconfig:"MINIO_USE_SSL" required:"true"`
	AccessKeyID      string `envconfig:"MINIO_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey  string `envconfig:"MINIO_SECRET_ACCESS_KEY" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load minio config: %w", err)
	}
	return &cfg, nil
}
