package customer

import (
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, data RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data LoginDto) (string, error)
}
