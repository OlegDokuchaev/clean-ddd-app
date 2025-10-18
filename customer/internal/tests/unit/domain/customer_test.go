package customer

import (
	"customer/internal/tests/testutils/mothers"
	"testing"
	"time"

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
		email       string
		phone       string
		password    string
		expectedErr error
	}{
		{
			name:     "Success",
			Cname:    "name",
			email:    "email@email.com",
			phone:    "+79032895555",
			password: "password",
		},
		{
			name:        "Failure: Empty name",
			Cname:       "",
			email:       "email@email.com",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerName,
		},
		{
			name:        "Failure: Name consists only of whitespaces",
			Cname:       "  ",
			email:       "email@email.com",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerName,
		},
		{
			name:        "Failure: Invalid phone",
			Cname:       "name",
			email:       "email@email.com",
			phone:       "invalid phone",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerPhone,
		},
		{
			name:        "Failure: Invalid email",
			Cname:       "name",
			email:       "invalid email",
			phone:       "+79032895555",
			password:    "password",
			expectedErr: customerDomain.ErrInvalidCustomerEmail,
		},
		{
			name:        "Failure: Empty password",
			email:       "email@email.com",
			Cname:       "name",
			phone:       "+79032895555",
			password:    "",
			expectedErr: customerDomain.ErrInvalidCustomerPassword,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.T().Parallel()

			customer, err := customerDomain.Create(tc.Cname, tc.phone, tc.email, tc.password)

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
				return mothers.DefaultCustomer()
			},
			newPassword:      "new password",
			expectedPassword: "new password",
			expectedErr:      nil,
		},
		{
			name: "Failure: Empty password",
			setup: func() *customerDomain.Customer {
				return mothers.CustomerWithPassword("password")
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
				return mothers.CustomerWithPassword("password")
			},
			password:    "password",
			expectedRes: true,
		},
		{
			name: "Failure: Invalid password",
			setup: func() *customerDomain.Customer {
				return mothers.CustomerWithPassword("password")
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

func (s *CustomerTestSuite) TestIsLocked() {
	now := time.Now()
	future := now.Add(30 * time.Minute)
	past := now.Add(-30 * time.Minute)

	tests := []struct {
		name        string
		setup       func() *customerDomain.Customer
		expectedRes bool
	}{
		{
			name: "LockedUntil in future -> true",
			setup: func() *customerDomain.Customer {
				return mothers.LockedUntilCustomer(future)
			},
			expectedRes: true,
		},
		{
			name: "LockedUntil in past -> false",
			setup: func() *customerDomain.Customer {
				return mothers.LockedUntilCustomer(past)
			},
			expectedRes: false,
		},
		{
			name: "No lock -> false",
			setup: func() *customerDomain.Customer {
				return mothers.DefaultCustomer()
			},
			expectedRes: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			c := tc.setup()

			result := c.IsLocked()

			require.Equal(s.T(), tc.expectedRes, result)
		})
	}
}

func (s *CustomerTestSuite) TestRegisterFailedAttempt() {
	lp := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: 1 * time.Hour}

	tests := []struct {
		name               string
		initialFailedCount int
		attempts           int
		expectFailedCount  int
		expectLocked       bool
	}{
		{
			name:               "First wrong attempt: count=1, not locked",
			initialFailedCount: 0,
			attempts:           1,
			expectFailedCount:  1,
			expectLocked:       false,
		},
		{
			name:               "Second wrong attempt: count=2, not locked",
			initialFailedCount: 1,
			attempts:           1,
			expectFailedCount:  2,
			expectLocked:       false,
		},
		{
			name:               "Reaches threshold: lock on 3rd",
			initialFailedCount: 2,
			attempts:           1,
			expectFailedCount:  3,
			expectLocked:       true,
		},
		{
			name:               "Already high count, another attempt keeps locked",
			initialFailedCount: 3,
			attempts:           1,
			expectFailedCount:  4,
			expectLocked:       true,
		},
		{
			name:               "Multiple attempts crossing threshold at once",
			initialFailedCount: 0,
			attempts:           3,
			expectFailedCount:  3,
			expectLocked:       true,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()

			c := mothers.DefaultCustomer()
			c.FailedCount = tc.initialFailedCount
			c.LockedUntil = nil

			for i := 0; i < tc.attempts; i++ {
				c.RegisterFailedAttempt(lp)
			}

			require.Equal(s.T(), tc.expectFailedCount, c.FailedCount)
			if tc.expectLocked {
				require.NotNil(s.T(), c.LockedUntil)
				require.True(s.T(), c.IsLocked())
			} else {
				require.True(s.T(), c.LockedUntil == nil || !c.IsLocked())
			}
		})
	}
}

func (s *CustomerTestSuite) TestResetFailedAttemptsAndUnlock() {
	now := time.Now()
	future := now.Add(45 * time.Minute)

	tests := []struct {
		name              string
		setup             func() *customerDomain.Customer
		expectFailedCount int
		expectLocked      bool
	}{
		{
			name: "Unlocks and resets counter from locked state",
			setup: func() *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				c.FailedCount = 5
				c.LockedUntil = &future
				return c
			},
			expectFailedCount: 0,
			expectLocked:      false,
		},
		{
			name: "No-op when already unlocked (just ensures zeros)",
			setup: func() *customerDomain.Customer {
				return mothers.DefaultCustomer()
			},
			expectFailedCount: 0,
			expectLocked:      false,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			c := tc.setup()

			c.ResetFailedAttempts()

			require.Equal(s.T(), tc.expectFailedCount, c.FailedCount)
			require.Equal(s.T(), tc.expectLocked, c.IsLocked())
			if !tc.expectLocked {
				require.Nil(s.T(), c.LockedUntil)
			}
		})
	}
}

func TestCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerTestSuite))
}
