package product

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type ImageUseCase interface {
	GetByID(ctx context.Context, productID uuid.UUID) (fileReader io.ReadCloser, contentType string, err error)
	UpdateByID(ctx context.Context, productID uuid.UUID, fileReader io.Reader, contentType string) error
}
