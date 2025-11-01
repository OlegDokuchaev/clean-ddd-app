package customer

import (
	customerDto "api-gateway/internal/domain/dtos/customer"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Register(ctx context.Context, data customerDto.RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data customerDto.LoginDto) (string, error)
	VerifyOtp(ctx context.Context, data customerDto.VerifyOtpDto) (string, error)
	RequestPasswordReset(ctx context.Context, email string) error
	CompletePasswordReset(ctx context.Context, token string, newPassword string) error
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
