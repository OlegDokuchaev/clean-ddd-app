package domain

import (
	courierDomain "courier/internal/domain/courier"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CourierDomainTestSuite struct {
	suite.Suite
}

func (s *CourierDomainTestSuite) TestSetPassword() {
	tests := []struct {
		name             string
		setup            func() *courierDomain.Courier
		newPassword      string
		expectedPassword string
		expectedErr      error
	}{
		{
			name: "Success",
			setup: func() *courierDomain.Courier {
				courier, err := courierDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return courier
			},
			newPassword:      "new password",
			expectedPassword: "new password",
			expectedErr:      nil,
		},
		{
			name: "Failure: Empty password",
			setup: func() *courierDomain.Courier {
				courier, err := courierDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return courier
			},
			newPassword:      "",
			expectedPassword: "password",
			expectedErr:      courierDomain.ErrInvalidCourierPassword,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			courier := tc.setup()

			err := courier.SetPassword(tc.newPassword)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
			}
			require.True(s.T(), courier.CheckPassword(tc.expectedPassword))
		})
	}
}

func (s *CourierDomainTestSuite) TestCheckPassword() {
	tests := []struct {
		name        string
		setup       func() *courierDomain.Courier
		password    string
		expectedRes bool
	}{
		{
			name: "Passwords are equal",
			setup: func() *courierDomain.Courier {
				courier, err := courierDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return courier
			},
			password:    "password",
			expectedRes: true,
		},
		{
			name: "Passwords are not equal",
			setup: func() *courierDomain.Courier {
				courier, err := courierDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return courier
			},
			password:    "new password",
			expectedRes: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			courier := tc.setup()

			res := courier.CheckPassword(tc.password)

			require.Equal(s.T(), tc.expectedRes, res)
		})
	}
}

func (s *CourierDomainTestSuite) TestCreate() {
	tests := []struct {
		name        string
		Cname       string
		phone       string
		password    string
		expectedErr error
	}{
		{
			name:        "Success",
			Cname:       "name",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: nil,
		},
		{
			name:        "Failure: Empty name",
			Cname:       "",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: courierDomain.ErrInvalidCourierName,
		},
		{
			name:        "Failure: Name consists only of whitespaces",
			Cname:       "  ",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: courierDomain.ErrInvalidCourierName,
		},
		{
			name:        "Failure: Invalid phone",
			Cname:       "name",
			phone:       "invalid phone",
			password:    "password",
			expectedErr: courierDomain.ErrInvalidCourierPhone,
		},
		{
			name:        "Failure: Invalid password",
			Cname:       "name",
			phone:       "+79032895555",
			password:    "",
			expectedErr: courierDomain.ErrInvalidCourierPassword,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()

			courier, err := courierDomain.Create(tc.Cname, tc.phone, tc.password)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), courier)
			}
		})
	}
}

func TestCourierDomainTestSuite(t *testing.T) {
	suite.Run(t, new(CourierDomainTestSuite))
}
