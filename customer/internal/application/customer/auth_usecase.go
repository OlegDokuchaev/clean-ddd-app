package customer

import (
	"context"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, data RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data LoginDto) (string, error)
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
