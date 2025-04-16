package product

import (
	"context"
	"github.com/google/uuid"
	"io"
)

type ImageUseCase interface {
	UpdateByID(ctx context.Context, productID uuid.UUID, fileReader io.Reader, contentType string) error
}
