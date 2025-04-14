package product

import (
	"context"
	"io"
	productApplication "warehouse/internal/application/product"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageServiceImpl struct {
	config *Config
	client *minio.Client
}

func NewImageService(config *Config, client *minio.Client) *ImageServiceImpl {
	return &ImageServiceImpl{
		config: config,
		client: client,
	}
}

func (s *ImageServiceImpl) Create(ctx context.Context) (string, error) {
	objectName := uuid.New().String()

	src := minio.CopySrcOptions{
		Bucket: s.config.BucketName,
		Object: s.config.DefaultImageName,
	}
	dst := minio.CopyDestOptions{
		Bucket: s.config.BucketName,
		Object: objectName,
	}
	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return "", parseError(err)
	}

	return objectName, nil
}

func (s *ImageServiceImpl) Update(ctx context.Context, path string, fileReader io.Reader, contentType string) error {
	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := s.client.PutObject(ctx, s.config.BucketName, path, fileReader, -1, options)
	if err != nil {
		return parseError(err)
	}

	return nil
}

var _ productApplication.ImageService = (*ImageServiceImpl)(nil)
