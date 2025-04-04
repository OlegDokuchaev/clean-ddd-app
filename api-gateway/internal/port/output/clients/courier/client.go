package courier

import (
	courierDto "api-gateway/internal/domain/dtos/courier"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Register(ctx context.Context, data courierDto.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data courierDto.LoginDto) (string, error)
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
