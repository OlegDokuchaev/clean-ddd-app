package product

import (
	"golang.org/x/net/context"
	"io"
)

type ImageService interface {
	Create(ctx context.Context) (string, error)
	Update(ctx context.Context, path string, fileReader io.Reader, contentType string) error
}
