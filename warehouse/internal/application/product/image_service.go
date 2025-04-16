package product

import (
	"io"

	"golang.org/x/net/context"
)

type ImageService interface {
	Create(ctx context.Context) (string, error)
	Get(ctx context.Context, path string) (fileReader io.ReadCloser, contentType string, err error)
	Update(ctx context.Context, path string, fileReader io.Reader, contentType string) error
}
