package courier

import (
	courierAuth "courier/internal/application/courier/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type TokenServiceMock struct {
	mock.Mock
}

func (m *TokenServiceMock) Generate(courierID uuid.UUID) (string, error) {
	args := m.Called(courierID)
	return args.String(0), args.Error(1)
}

func (m *TokenServiceMock) Decode(token string) (uuid.UUID, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

var _ courierAuth.TokenService = (*TokenServiceMock)(nil)
