package testutils

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
	minioContainer "github.com/testcontainers/testcontainers-go/modules/minio"
)

const (
	minioImage     = "minio/minio:RELEASE.2024-01-16T16-07-38Z"
	minioAccessKey = "minioadmin"
	minioSecretKey = "minioadmin"
)

type TestMinio struct {
	Client    *minio.Client
	container testcontainers.Container
}

func (d *TestMinio) Close(ctx context.Context) error {
	return d.container.Terminate(ctx)
}

func setupMinioContainer(ctx context.Context) (testcontainers.Container, error) {
	return minioContainer.Run(
		ctx,
		minioImage,
		minioContainer.WithUsername(minioAccessKey),
		minioContainer.WithPassword(minioSecretKey),
	)
}

func createMinioEndpoint(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := container.MappedPort(ctx, "9000/tcp")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

func createMinioClient(ctx context.Context, container testcontainers.Container) (*minio.Client, error) {
	endpoint, err := createMinioEndpoint(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port for MinIO: %w", err)
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return minioClient, nil
}

func NewTestMinio(ctx context.Context) (*TestMinio, error) {
	container, err := setupMinioContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup minio: %w", err)
	}

	client, err := createMinioClient(ctx, container)
	if err != nil {
		if err := container.Terminate(ctx); err != nil {
			return nil, fmt.Errorf("failed to terminate minio: %w", err)
		}
		return nil, err
	}

	return &TestMinio{
		Client:    client,
		container: container,
	}, nil
}

func GetFileHash(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", fmt.Errorf("failed to copy data to hash: %w", err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
