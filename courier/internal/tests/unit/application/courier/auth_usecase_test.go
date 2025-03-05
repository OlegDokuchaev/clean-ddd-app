package courier

import (
	"context"
	courierAuth "courier/internal/application/courier/auth"
	courierDomain "courier/internal/domain/courier"
	courierMock "courier/internal/mocks/courier"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *AuthUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *AuthUseCaseTestSuite) createTestCourier() *courierDomain.Courier {
	courier, err := courierDomain.Create("test", "+79032895555", "password")
	require.NoError(s.T(), err)
	return courier
}

func (s *AuthUseCaseTestSuite) TestRegister() {
	tests := []struct {
		name        string
		data        courierAuth.RegisterDto
		setup       func(repo *courierMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: courierAuth.RegisterDto{
				Name:     "test",
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *courierMock.RepositoryMock) {
				repo.On("Create", s.ctx, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Create courier error",
			data: courierAuth.RegisterDto{
				Name:     "test",
				Phone:    "invalid phone",
				Password: "password",
			},
			setup:       func(repo *courierMock.RepositoryMock) {},
			expectedErr: courierDomain.ErrInvalidCourierPhone,
		},
		{
			name: "Failure: Create courier error",
			data: courierAuth.RegisterDto{
				Name:     "test",
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *courierMock.RepositoryMock) {
				repo.On("Create", s.ctx, mock.Anything).Return(errors.New("create courier error"))
			},
			expectedErr: errors.New("create courier error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(courierMock.RepositoryMock)
			token := new(courierMock.TokenServiceMock)
			uc := courierAuth.NewUseCase(repo, token)
			tc.setup(repo)

			courierID, err := uc.Register(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEqual(s.T(), uuid.Nil, courierID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Equal(s.T(), uuid.Nil, courierID)
			}

			repo.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestLogin() {
	tests := []struct {
		name        string
		data        courierAuth.LoginDto
		setup       func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: courierAuth.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) {
				courier := s.createTestCourier()
				repo.On("GetByPhone", s.ctx, courier.Phone).Return(courier, nil)
				token.On("Generate", courier.ID).Return("token", nil)
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Get courier by phone error",
			data: courierAuth.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) {
				repo.On("GetByPhone", s.ctx, "+79032895555").Return((*courierDomain.Courier)(nil), errors.New("get courier by phone error"))
			},
			expectedErr: errors.New("get courier by phone error"),
		},
		{
			name: "Failure: Invalid courier password",
			data: courierAuth.LoginDto{
				Phone:    "+79032895555",
				Password: "invalid password",
			},
			setup: func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) {
				courier := s.createTestCourier()
				repo.On("GetByPhone", s.ctx, courier.Phone).Return(courier, nil)
			},
			expectedErr: courierDomain.ErrInvalidCourierPassword,
		},
		{
			name: "Failure: Generate token error",
			data: courierAuth.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) {
				courier := s.createTestCourier()
				repo.On("GetByPhone", s.ctx, courier.Phone).Return(courier, nil)
				token.On("Generate", courier.ID).Return("", errors.New("generate token error"))
			},
			expectedErr: errors.New("generate token error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(courierMock.RepositoryMock)
			tokenService := new(courierMock.TokenServiceMock)
			uc := courierAuth.NewUseCase(repo, tokenService)
			tc.setup(repo, tokenService)

			token, err := uc.Login(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEmpty(s.T(), token)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Empty(s.T(), token)
			}

			repo.AssertExpectations(s.T())
			tokenService.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestAuthenticate() {
	tests := []struct {
		name        string
		data        string
		setup       func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) uuid.UUID
		expectedErr error
	}{
		{
			name: "Success",
			data: "token",
			setup: func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) uuid.UUID {
				courierID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				token.On("Decode", "token").Return(courierID, nil)
				return courierID
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Decode token error",
			data: "token",
			setup: func(repo *courierMock.RepositoryMock, token *courierMock.TokenServiceMock) uuid.UUID {
				token.On("Decode", "token").Return(uuid.Nil, errors.New("decode token error"))
				return uuid.Nil
			},
			expectedErr: errors.New("decode token error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(courierMock.RepositoryMock)
			tokenService := new(courierMock.TokenServiceMock)
			uc := courierAuth.NewUseCase(repo, tokenService)
			expectedCourierID := tc.setup(repo, tokenService)

			courierID, err := uc.Authenticate(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Equal(s.T(), courierID, expectedCourierID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Equal(s.T(), uuid.Nil, courierID)
			}

			repo.AssertExpectations(s.T())
			tokenService.AssertExpectations(s.T())
		})
	}
}

func TestAuthUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}
