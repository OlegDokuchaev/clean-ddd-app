package customer

import (
	"context"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, data RegisterDto) (uuid.UUID, error)
	Login(ctx context.Context, data LoginDto) (string, error)
	VerifyOtp(ctx context.Context, data VerifyOtpDto) (string, error)
	RequestPasswordReset(ctx context.Context, email string) error
	CompletePasswordReset(ctx context.Context, token string, newPassword string) error
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}
