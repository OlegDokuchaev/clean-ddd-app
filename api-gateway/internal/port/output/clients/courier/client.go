package courier

import (
	"api-gateway/internal/domain/dtos"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Register(ctx context.Context, data dtos.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data dtos.LoginDto) (string, error)
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
