package auth

import (
	"customer/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TokenManagerTestSuite struct {
	suite.Suite

	customerID uuid.UUID
}

func (s *TokenManagerTestSuite) SetupSuite() {
	s.customerID = uuid.New()
}

func (s *TokenManagerTestSuite) TestGenerate() {
	tests := []struct {
		name        string
		config      *auth.Config
		customerID  uuid.UUID
		expectedErr error
	}{
		{
			name: "Success",
			config: &auth.Config{
				SigningKey: "test_signing_key",
				TokenTTL:   time.Hour,
			},
			customerID:  s.customerID,
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			tokenManager := auth.NewTokenManager(tc.config)

			token, err := tokenManager.Generate(tc.customerID)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			} else {
				require.NoError(s.T(), err)
				require.NotEmpty(s.T(), token)

				returnedCustomerID, err := tokenManager.Decode(token)
				require.NoError(s.T(), err)
				assert.Equal(s.T(), tc.customerID, returnedCustomerID)
			}
		})
	}
}

func (s *TokenManagerTestSuite) TestDecode() {
	tests := []struct {
		name        string
		setup       func() (*auth.TokenManagerImpl, string)
		expectedID  uuid.UUID
		expectedErr error
	}{
		{
			name: "Success: Valid token",
			setup: func() (*auth.TokenManagerImpl, string) {
				config := &auth.Config{
					SigningKey: "test_signing_key",
					TokenTTL:   time.Hour,
				}
				tokenManager := auth.NewTokenManager(config)
				token, err := tokenManager.Generate(s.customerID)
				require.NoError(s.T(), err)
				return tokenManager, token
			},
			expectedID:  s.customerID,
			expectedErr: nil,
		},
		{
			name: "Error: Invalid token format",
			setup: func() (*auth.TokenManagerImpl, string) {
				config := &auth.Config{
					SigningKey: "test_signing_key",
					TokenTTL:   time.Hour,
				}
				tokenManager := auth.NewTokenManager(config)
				return tokenManager, "invalid.token.string"
			},
			expectedID:  uuid.Nil,
			expectedErr: auth.ErrInvalidToken,
		},
		{
			name: "Error: Expired token",
			setup: func() (*auth.TokenManagerImpl, string) {
				config := &auth.Config{
					SigningKey: "test_signing_key",
					TokenTTL:   -time.Hour,
				}
				tokenManager := auth.NewTokenManager(config)
				token, err := tokenManager.Generate(s.customerID)
				require.NoError(s.T(), err)
				return tokenManager, token
			},
			expectedID:  uuid.Nil,
			expectedErr: auth.ErrInvalidToken,
		},
		{
			name: "Error - Invalid signing method",
			setup: func() (*auth.TokenManagerImpl, string) {
				config := &auth.Config{
					SigningKey: "test_signing_key",
					TokenTTL:   time.Hour,
				}
				tokenManager := auth.NewTokenManager(config)

				tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb3VyaWVyX2lkIjoiOTZkMzE0ZWUtNWY3MC00NjYzLWIxYTktN2U2MzY2YThjZmRiIiwiZXhwIjoxNzE4OTE2MzY5LCJpYXQiOjE3MTY1NTQ3Njl9.invalid_signature"

				return tokenManager, tokenString
			},
			expectedID:  uuid.Nil,
			expectedErr: auth.ErrInvalidToken,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			tokenManager, token := tc.setup()

			customerID, err := tokenManager.Decode(token)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), tc.expectedID, customerID)
			}
		})
	}
}

func TestTokenManagerTestSuite(t *testing.T) {
	suite.Run(t, new(TokenManagerTestSuite))
}
