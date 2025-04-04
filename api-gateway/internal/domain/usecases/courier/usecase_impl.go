package courier

import (
	courierDto "api-gateway/internal/domain/dtos/courier"
	courierClient "api-gateway/internal/port/output/clients/courier"
	"context"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	courierClient courierClient.Client
}

func NewUseCase(
	courierClient courierClient.Client,
) UseCase {
	return &UseCaseImpl{
		courierClient: courierClient,
	}
}

func (u *UseCaseImpl) Register(ctx context.Context, data courierDto.RegisterDto) (uuid.UUID, error) {
	courierID, err := u.courierClient.Register(ctx, data)
	if err != nil {
		return uuid.Nil, err
	}

	return courierID, nil
}

func (u *UseCaseImpl) Login(ctx context.Context, data courierDto.LoginDto) (string, error) {
	token, err := u.courierClient.Login(ctx, data)
	if err != nil {
		return "", err
	}

	return token, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
