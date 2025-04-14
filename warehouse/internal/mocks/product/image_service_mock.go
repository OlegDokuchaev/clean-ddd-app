package product

import (
	"context"
	"io"
	productApplication "warehouse/internal/application/product"

	"github.com/stretchr/testify/mock"
)

type ImageServiceMock struct {
	mock.Mock
}

func (m *ImageServiceMock) Create(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *ImageServiceMock) Update(ctx context.Context, path string, fileReader io.Reader, contentType string) error {
	args := m.Called(ctx, path, fileReader, contentType)
	return args.Error(0)
}

var _ productApplication.ImageService = (*ImageServiceMock)(nil)
