//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/db/migrations"
	customerRepository "customer/internal/infrastructure/repository/customer"
	"customer/internal/tests/testutils"
	"customer/internal/tests/testutils/mothers"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *CustomerRepositoryTestSuite) SetupSuite() {
	cfg, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	s.testDB, err = testutils.NewTestDB(s.ctx, cfg)
	require.NoError(s.T(), err)
}

func (s *CustomerRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *CustomerRepositoryTestSuite) getRepo() customerDomain.Repository {
	return customerRepository.New(s.testDB.DB)
}

func (s *CustomerRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		setup         func(repo customerDomain.Repository) *customerDomain.Customer
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				return mothers.DefaultCustomer()
			},
			expectedError: nil,
		},
		{
			name: "Failure: Duplicate primary key (same ID)",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c))
				return c
			},
			expectedError: customerRepository.ErrCustomerAlreadyExists,
		},
		{
			name: "Failure: Phone already exists",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c1 := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c1))
				return mothers.CustomerWithPhone(c1.Phone)
			},
			expectedError: customerRepository.ErrCustomerPhoneAlreadyExists,
		},
		{
			name: "Failure: Email already exists",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c1 := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c1))
				return mothers.CustomerWithEmail(c1.Email)
			},
			expectedError: customerRepository.ErrCustomerEmailAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			c := tc.setup(repo)

			err := repo.Create(s.ctx, c)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				created, err := repo.GetByPhone(s.ctx, c.Phone)
				require.NoError(s.T(), err)
				require.Equal(s.T(), c.ID, created.ID)
				require.Equal(s.T(), c.Email, created.Email)
			}
		})
	}
}

func (s *CustomerRepositoryTestSuite) TestGetByPhone() {
	tests := []struct {
		name          string
		setup         func(repo customerDomain.Repository) *customerDomain.Customer
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c))
				return c
			},
			expectedError: nil,
		},
		{
			name: "Failure: Customer not found",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				return mothers.DefaultCustomer()
			},
			expectedError: customerRepository.ErrCustomerNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			created := tc.setup(repo)

			got, err := repo.GetByPhone(s.ctx, created.Phone)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
				require.Nil(s.T(), got)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), created.ID, got.ID)
			}
		})
	}
}

func (s *CustomerRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name          string
		setup         func(repo customerDomain.Repository) *customerDomain.Customer
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c))
				return c
			},
			expectedError: nil,
		},
		{
			name: "Failure: Customer not found",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				return mothers.DefaultCustomer()
			},
			expectedError: customerRepository.ErrCustomerNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			created := tc.setup(repo)

			got, err := repo.GetByID(s.ctx, created.ID)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
				require.Nil(s.T(), got)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), created.ID, got.ID)
			}
		})
	}
}

func (s *CustomerRepositoryTestSuite) TestGetByEmail() {
	tests := []struct {
		name          string
		setup         func(repo customerDomain.Repository) *customerDomain.Customer
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c))
				return c
			},
			expectedError: nil,
		},
		{
			name: "Failure: Customer not found",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				return mothers.DefaultCustomer()
			},
			expectedError: customerRepository.ErrCustomerNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			created := tc.setup(repo)

			got, err := repo.GetByEmail(s.ctx, created.Email)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
				require.Nil(s.T(), got)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), created.ID, got.ID)
			}
		})
	}
}

func (s *CustomerRepositoryTestSuite) TestSave() {
	now := time.Now()
	lockFuture := now.Add(45 * time.Minute)

	tests := []struct {
		name          string
		setup         func(repo customerDomain.Repository) *customerDomain.Customer
		mutate        func(c *customerDomain.Customer)
		verify        func(repo customerDomain.Repository, c *customerDomain.Customer)
		expectedError error
	}{
		{
			name: "Success: update counters and lock",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c))
				return c
			},
			mutate: func(c *customerDomain.Customer) {
				c.FailedCount = 5
				c.LockedUntil = &lockFuture
				c.MustChangePassword = true
			},
			verify: func(repo customerDomain.Repository, c *customerDomain.Customer) {
				got, err := repo.GetByID(s.ctx, c.ID)
				require.NoError(s.T(), err)
				require.Equal(s.T(), 5, got.FailedCount)
				require.NotNil(s.T(), got.LockedUntil)
				require.True(s.T(), got.MustChangePassword)
			},
			expectedError: nil,
		},
		{
			name: "Success: clear lock (LockedUntil -> NULL) and reset counter to zero",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c := mothers.DefaultCustomer()
				c.FailedCount = 3
				t := now.Add(10 * time.Minute)
				c.LockedUntil = &t
				c.MustChangePassword = true
				require.NoError(s.T(), repo.Create(s.ctx, c))
				return c
			},
			mutate: func(c *customerDomain.Customer) {
				c.FailedCount = 0
				c.LockedUntil = nil
				c.MustChangePassword = false
			},
			verify: func(repo customerDomain.Repository, c *customerDomain.Customer) {
				got, err := repo.GetByID(s.ctx, c.ID)
				require.NoError(s.T(), err)
				require.Equal(s.T(), 0, got.FailedCount)
				require.Nil(s.T(), got.LockedUntil)
				require.False(s.T(), got.MustChangePassword)
			},
			expectedError: nil,
		},
		{
			name: "Failure: change email to existing -> unique error",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c1 := mothers.DefaultCustomer()
				c2 := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c1))
				require.NoError(s.T(), repo.Create(s.ctx, c2))
				c2.Email = c1.Email
				return c2
			},
			mutate: func(c *customerDomain.Customer) {
				c.Name = c.Name + " updated"
			},
			verify:        func(repo customerDomain.Repository, c *customerDomain.Customer) {}, // не дойдём
			expectedError: customerRepository.ErrCustomerEmailAlreadyExists,
		},
		{
			name: "Failure: change phone to existing -> unique error",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				c1 := mothers.DefaultCustomer()
				c2 := mothers.DefaultCustomer()
				require.NoError(s.T(), repo.Create(s.ctx, c1))
				require.NoError(s.T(), repo.Create(s.ctx, c2))
				c2.Phone = c1.Phone
				return c2
			},
			mutate:        func(c *customerDomain.Customer) {},
			verify:        func(repo customerDomain.Repository, c *customerDomain.Customer) {},
			expectedError: customerRepository.ErrCustomerPhoneAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			c := tc.setup(repo)
			tc.mutate(c)

			err := repo.Save(s.ctx, c)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
				return
			}

			require.NoError(s.T(), err)
			tc.verify(repo, c)
		})
	}
}

func TestCustomerRepository(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
