package customer

import (
	"testing"

	customerDomain "customer/internal/domain/customer"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CustomerTestSuite struct {
	suite.Suite
}

func (s *CustomerTestSuite) TestCreateCustomer() {
	tests := []struct {
		name        string
		Cname       string
		phone       string
		password    string
		expectedErr error
	}{
		{
			name:     "Success",
			Cname:    "name",
			phone:    "+79032895555",
			password: "password",
		},
		{
			name:        "Failure: Empty name",
			Cname:       "",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerName,
		},
		{
			name:        "Failure: Name consists only of whitespaces",
			Cname:       "  ",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerName,
		},
		{
			name:        "Failure: Invalid phone",
			Cname:       "name",
			phone:       "invalid phone",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerPhone,
		},
		{
			name:        "Failure: Empty password",
			Cname:       "name",
			phone:       "+79032895555",
			password:    "",
			expectedErr: customerDomain.ErrInvalidCustomerPassword,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.T().Parallel()

			customer, err := customerDomain.Create(tc.Cname, tc.phone, tc.password)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), customer)
			}
		})
	}
}

func (s *CustomerTestSuite) TestSetPassword() {
	tests := []struct {
		name             string
		setup            func() *customerDomain.Customer
		newPassword      string
		expectedPassword string
		expectedErr      error
	}{
		{
			name: "Success",
			setup: func() *customerDomain.Customer {
				customer, err := customerDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return customer
			},
			newPassword:      "new password",
			expectedPassword: "new password",
			expectedErr:      nil,
		},
		{
			name: "Failure: Empty password",
			setup: func() *customerDomain.Customer {
				customer, err := customerDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return customer
			},
			newPassword:      "",
			expectedPassword: "password",
			expectedErr:      customerDomain.ErrInvalidCustomerPassword,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			customer := tc.setup()

			err := customer.SetPassword(tc.newPassword)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
			}
			require.True(s.T(), customer.CheckPassword(tc.expectedPassword))
		})
	}
}

func (s *CustomerTestSuite) TestCheckPassword() {
	tests := []struct {
		name        string
		setup       func() *customerDomain.Customer
		password    string
		expectedRes bool
	}{
		{
			name: "Success",
			setup: func() *customerDomain.Customer {
				customer, err := customerDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return customer
			},
			password:    "password",
			expectedRes: true,
		},
		{
			name: "Failure: Invalid password",
			setup: func() *customerDomain.Customer {
				customer, err := customerDomain.Create("test", "+79032895555", "password")
				require.NoError(s.T(), err)
				return customer
			},
			password:    "wrong password",
			expectedRes: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			customer := tc.setup()

			result := customer.CheckPassword(tc.password)

			require.Equal(s.T(), tc.expectedRes, result)
		})
	}
}

func TestCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerTestSuite))
}
