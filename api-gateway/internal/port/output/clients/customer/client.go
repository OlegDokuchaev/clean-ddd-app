package customer

import (
	customerUseCase "api-gateway/internal/domain/usecases/customer"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Register(ctx context.Context, data customerUseCase.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data customerUseCase.LoginDto) (string, error)
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
