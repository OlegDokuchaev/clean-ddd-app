package customer

import (
	"context"
	customerDomain "customer/internal/domain/customer"

	"github.com/google/uuid"
)

type AuthUseCaseImpl struct {
	repo         customerDomain.Repository
	tokenManager TokenManager
}

func NewAuthUseCase(repo customerDomain.Repository, tokenManager TokenManager) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (u *AuthUseCaseImpl) Register(ctx context.Context, data RegisterDto) (uuid.UUID, error) {
	customer, err := customerDomain.Create(data.Name, data.Phone, data.Password)
	if err != nil {
		return uuid.Nil, err
	}

	if err = u.repo.Create(ctx, customer); err != nil {
		return uuid.Nil, err
	}

	return customer.ID, nil
}

func (u *AuthUseCaseImpl) Login(ctx context.Context, data LoginDto) (string, error) {
	customer, err := u.repo.GetByPhone(ctx, data.Phone)
	if err != nil {
		return "", err
	}

	if valid := customer.CheckPassword(data.Password); !valid {
		return "", customerDomain.ErrInvalidCustomerPassword
	}

	return u.tokenManager.Generate(customer.ID)
}

func (u *AuthUseCaseImpl) Authenticate(_ context.Context, token string) (uuid.UUID, error) {
	return u.tokenManager.Decode(token)
}

var _ AuthUseCase = (*AuthUseCaseImpl)(nil)
