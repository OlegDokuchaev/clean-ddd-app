//go:build integration

package repository

import (
	"context"
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/db/migrations"
	customerRepository "customer/internal/infrastructure/repository/customer"
	"customer/internal/tests/testutils"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *CustomerRepositoryTestSuite) SetupSuite() {
	config, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	s.testDB, err = testutils.NewTestDB(s.ctx, config)
	require.NoError(s.T(), err)
}

func (s *CustomerRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Container.Terminate(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *CustomerRepositoryTestSuite) getRepo() customerDomain.Repository {
	return customerRepository.New(s.testDB.DB)
}

func (s *CustomerRepositoryTestSuite) createRandomPhone() string {
	phone := "+"
	for i := 0; i < 10; i++ {
		phone += fmt.Sprintf("%d", rand.Intn(10))
	}
	return phone
}

func (s *CustomerRepositoryTestSuite) createTestCustomer(phone string) *customerDomain.Customer {
	customer, err := customerDomain.Create("John Doe", phone, "john.doe@example.com")
	require.NoError(s.T(), err)
	return customer
}

func (s *CustomerRepositoryTestSuite) createTestCustomerInDb(phone string, repo customerDomain.Repository) *customerDomain.Customer {
	customer := s.createTestCustomer(phone)
	err := repo.Create(s.ctx, customer)
	require.NoError(s.T(), err)
	return customer
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
				return s.createTestCustomer(s.createRandomPhone())
			},
			expectedError: nil,
		},
		{
			name: "Failure: Duplicate customer",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				return s.createTestCustomerInDb(s.createRandomPhone(), repo)
			},
			expectedError: customerRepository.ErrCustomerAlreadyExists,
		},
		{
			name: "Failure: Phone number already exists",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				phone := s.createRandomPhone()

				customer := s.createTestCustomer(phone)
				err := repo.Create(s.ctx, customer)
				require.NoError(s.T(), err)

				return s.createTestCustomer(phone)
			},
			expectedError: customerRepository.ErrCustomerPhoneAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			customer := test.setup(repo)

			err := repo.Create(s.ctx, customer)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				createdCustomer, err := repo.GetByPhone(s.ctx, customer.Phone)
				require.NoError(s.T(), err)
				require.Equal(s.T(), customer.ID, createdCustomer.ID)
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
				return s.createTestCustomerInDb(s.createRandomPhone(), repo)
			},
			expectedError: nil,
		},
		{
			name: "Failure: Customer not found",
			setup: func(repo customerDomain.Repository) *customerDomain.Customer {
				return s.createTestCustomer(s.createRandomPhone())
			},
			expectedError: customerRepository.ErrCustomerNotFound,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			customer := test.setup(repo)

			createdCustomer, err := repo.GetByPhone(s.ctx, customer.Phone)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), customer.ID, createdCustomer.ID)
			}
		})
	}
}

func TestCustomerRepository(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
