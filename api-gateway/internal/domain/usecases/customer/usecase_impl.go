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

func NewUseCase(
	customerClient customerClient.Client,
) UseCase {
	return &UseCaseImpl{
		customerClient: customerClient,
	}
}

func (u *UseCaseImpl) Register(ctx context.Context, data customerDto.RegisterDto) (uuid.UUID, error) {
	customerID, err := u.customerClient.Register(ctx, data)
	if err != nil {
		return uuid.Nil, err
	}

	return customerID, nil
}

func (u *UseCaseImpl) Login(ctx context.Context, data customerDto.LoginDto) (string, error) {
	token, err := u.customerClient.Login(ctx, data)
	if err != nil {
		return "", err
	}

	return token, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
