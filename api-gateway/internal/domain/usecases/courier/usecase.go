package courier

import (
	"api-gateway/internal/domain/dtos"
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, data dtos.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data dtos.LoginDto) (string, error)
}
