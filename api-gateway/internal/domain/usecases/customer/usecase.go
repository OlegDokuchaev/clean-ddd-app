package customer

import (
	customerDto "api-gateway/internal/domain/dtos/customer"
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, data customerDto.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data customerDto.LoginDto) (string, error)
}
