package customer

import (
	customerDto "api-gateway/internal/domain/dtos/customer"
	customerClient "api-gateway/internal/port/output/clients/customer"
	"context"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	customerClient customerClient.Client
}

func NewUseCase(customerClient customerClient.Client) UseCase {
	return &UseCaseImpl{customerClient: customerClient}
}

func (u *UseCaseImpl) Register(ctx context.Context, data customerDto.RegisterDto) (uuid.UUID, error) {
	return u.customerClient.Register(ctx, data)
}

func (u *UseCaseImpl) Login(ctx context.Context, data customerDto.LoginDto) (string, error) {
	return u.customerClient.Login(ctx, data)
}

func (u *UseCaseImpl) VerifyOtp(ctx context.Context, data customerDto.VerifyOtpDto) (string, error) {
	return u.customerClient.VerifyOtp(ctx, data)
}

func (u *UseCaseImpl) RequestPasswordReset(ctx context.Context, email string) error {
	return u.customerClient.RequestPasswordReset(ctx, email)
}

func (u *UseCaseImpl) CompletePasswordReset(ctx context.Context, token string, newPassword string) error {
	return u.customerClient.CompletePasswordReset(ctx, token, newPassword)
}

var _ UseCase = (*UseCaseImpl)(nil)
