package courier

import (
	courierDto "api-gateway/internal/domain/dtos/courier"
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, data courierDto.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data courierDto.LoginDto) (string, error)
}
