package courier

import (
	courierUseCase "api-gateway/internal/domain/usecases/courier"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Register(ctx context.Context, data courierUseCase.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data courierUseCase.LoginDto) (string, error)
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
