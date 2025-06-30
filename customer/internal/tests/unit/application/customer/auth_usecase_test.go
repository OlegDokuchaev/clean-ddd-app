package customer

import (
	"context"
	customerApplication "customer/internal/application/customer"
	customerDomain "customer/internal/domain/customer"
	customerMock "customer/internal/mocks/customer"
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

func (s *AuthUseCaseTestSuite) createTestCustomer() *customerDomain.Customer {
	customer, err := customerDomain.Create("test", "+79032895555", "password")
	require.NoError(s.T(), err)
	return customer
}

func (s *AuthUseCaseTestSuite) TestRegister() {
	tests := []struct {
		name        string
		data        customerApplication.RegisterDto
		setup       func(repo *customerMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: customerApplication.RegisterDto{
				Name:     "test",
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock) {
				repo.On("Create", s.ctx, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Create customer error",
			data: customerApplication.RegisterDto{
				Name:     "test",
				Phone:    "invalid phone",
				Password: "password",
			},
			setup:       func(repo *customerMock.RepositoryMock) {},
			expectedErr: customerDomain.ErrInvalidCustomerPhone,
		},
		{
			name: "Failure: Repository create customer error",
			data: customerApplication.RegisterDto{
				Name:     "test",
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock) {
				repo.On("Create", s.ctx, mock.Anything).
					Return(errors.New("create customer error"))
			},
			expectedErr: errors.New("create customer error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			uc := customerApplication.NewAuthUseCase(repo, tokenManager)
			tc.setup(repo)

			customerID, err := uc.Register(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEqual(s.T(), uuid.Nil, customerID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestLogin() {
	tests := []struct {
		name        string
		data        customerApplication.LoginDto
		setup       func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				customer := s.createTestCustomer()
				repo.On("GetByPhone", s.ctx, customer.Phone).Return(customer, nil)
				token.On("Generate", customer.ID).Return("token", nil)
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Get customer by phone error",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				repo.On("GetByPhone", s.ctx, "+79032895555").
					Return((*customerDomain.Customer)(nil), errors.New("get customer by phone error"))
			},
			expectedErr: errors.New("get customer by phone error"),
		},
		{
			name: "Failure: Invalid customer password",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "invalid password",
			},
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				customer := s.createTestCustomer()
				repo.On("GetByPhone", s.ctx, customer.Phone).Return(customer, nil)
			},
			expectedErr: customerDomain.ErrInvalidCustomerPassword,
		},
		{
			name: "Failure: Generate token error",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				customer := s.createTestCustomer()
				repo.On("GetByPhone", s.ctx, customer.Phone).Return(customer, nil)
				token.On("Generate", customer.ID).
					Return("", errors.New("generate token error"))
			},
			expectedErr: errors.New("generate token error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			uc := customerApplication.NewAuthUseCase(repo, tokenManager)
			tc.setup(repo, tokenManager)

			token, err := uc.Login(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEmpty(s.T(), token)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			tokenManager.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestAuthenticate() {
	tests := []struct {
		name        string
		data        string
		setup       func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) uuid.UUID
		expectedErr error
	}{
		{
			name: "Success",
			data: "token",
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) uuid.UUID {
				customerID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				token.On("Decode", "token").Return(customerID, nil)
				return customerID
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Decode token error",
			data: "token",
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) uuid.UUID {
				token.On("Decode", "token").
					Return(uuid.Nil, errors.New("decode token error"))
				return uuid.Nil
			},
			expectedErr: errors.New("decode token error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			uc := customerApplication.NewAuthUseCase(repo, tokenManager)
			expectedCustomerID := tc.setup(repo, tokenManager)

			customerID, err := uc.Authenticate(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Equal(s.T(), customerID, expectedCustomerID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			tokenManager.AssertExpectations(s.T())
		})
	}
}

func TestAuthUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}
