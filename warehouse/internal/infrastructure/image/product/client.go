package product

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(config *Config) (*minio.Client, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}
