package auth

import (
	"context"
	courierDomain "courier/internal/domain/courier"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	repo         courierDomain.Repository
	tokenService TokenService
}

func NewUseCase(repo courierDomain.Repository, tokenService TokenService) *UseCaseImpl {
	return &UseCaseImpl{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (u *UseCaseImpl) Register(ctx context.Context, data RegisterDto) (uuid.UUID, error) {
	courier, err := courierDomain.Create(data.Name, data.Phone, data.Password)
	if err != nil {
		return uuid.Nil, err
	}

	if err = u.repo.Create(ctx, courier); err != nil {
		return uuid.Nil, err
	}

	return courier.ID, nil
}

func (u *UseCaseImpl) Login(ctx context.Context, data LoginDto) (string, error) {
	courier, err := u.repo.GetByPhone(ctx, data.Phone)
	if err != nil {
		return "", err
	}

	if valid := courier.CheckPassword(data.Password); !valid {
		return "", courierDomain.ErrInvalidCourierPassword
	}

	token, err := u.tokenService.Generate(courier.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UseCaseImpl) Authenticate(_ context.Context, token string) (uuid.UUID, error) {
	courierID, err := u.tokenService.Decode(token)
	if err != nil {
		return uuid.Nil, err
	}
	return courierID, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
