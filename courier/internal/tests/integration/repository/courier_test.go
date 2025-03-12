//go:build integration

package repository

import (
	"context"
	courierDomain "courier/internal/domain/courier"
	"courier/internal/infrastructure/db/migrations"
	courierRepository "courier/internal/infrastructure/repository/courier"
	"courier/internal/tests/testutils"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CourierRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *CourierRepositoryTestSuite) SetupSuite() {
	config, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	s.testDB, err = testutils.NewTestDB(s.ctx, config)
	require.NoError(s.T(), err)
}

func (s *CourierRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Container.Terminate(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *CourierRepositoryTestSuite) getRepo() courierDomain.Repository {
	return courierRepository.New(s.testDB.DB)
}

func (s *CourierRepositoryTestSuite) createRandomPhone() string {
	phone := "+"
	for i := 0; i < 10; i++ {
		phone += fmt.Sprintf("%d", rand.Intn(10))
	}
	return phone
}

func (s *CourierRepositoryTestSuite) createTestCourier(phone string) *courierDomain.Courier {
	courier, err := courierDomain.Create("John Doe", phone, "john.doe@example.com")
	require.NoError(s.T(), err)
	return courier
}

func (s *CourierRepositoryTestSuite) createTestCourierInDb(phone string, repo courierDomain.Repository) *courierDomain.Courier {
	courier := s.createTestCourier(phone)
	err := repo.Create(s.ctx, courier)
	require.NoError(s.T(), err)
	return courier
}

func (s *CourierRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		setup         func(repo courierDomain.Repository) *courierDomain.Courier
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo courierDomain.Repository) *courierDomain.Courier {
				return s.createTestCourier(s.createRandomPhone())
			},
			expectedError: nil,
		},
		{
			name: "Failure: Duplicate courier",
			setup: func(repo courierDomain.Repository) *courierDomain.Courier {
				return s.createTestCourierInDb(s.createRandomPhone(), repo)
			},
			expectedError: courierRepository.ErrCourierAlreadyExists,
		},
		{
			name: "Failure: Phone number already exists",
			setup: func(repo courierDomain.Repository) *courierDomain.Courier {
				phone := s.createRandomPhone()

				courier := s.createTestCourier(phone)
				err := repo.Create(s.ctx, courier)
				require.NoError(s.T(), err)

				return s.createTestCourier(phone)
			},
			expectedError: courierRepository.ErrCourierPhoneAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			courier := test.setup(repo)

			err := repo.Create(s.ctx, courier)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				createdCourier, err := repo.GetByPhone(s.ctx, courier.Phone)
				require.NoError(s.T(), err)
				require.Equal(s.T(), courier.ID, createdCourier.ID)
			}
		})
	}
}

func (s *CourierRepositoryTestSuite) TestGetByPhone() {
	tests := []struct {
		name          string
		setup         func(repo courierDomain.Repository) *courierDomain.Courier
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo courierDomain.Repository) *courierDomain.Courier {
				return s.createTestCourierInDb(s.createRandomPhone(), repo)
			},
			expectedError: nil,
		},
		{
			name: "Failure: Courier not found",
			setup: func(repo courierDomain.Repository) *courierDomain.Courier {
				return s.createTestCourier(s.createRandomPhone())
			},
			expectedError: courierRepository.ErrCourierNotFound,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			courier := test.setup(repo)

			createdCourier, err := repo.GetByPhone(s.ctx, courier.Phone)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), courier.ID, createdCourier.ID)
			}
		})
	}
}

func (s *CourierRepositoryTestSuite) TestGetAll() {
	tests := []struct {
		name          string
		setup         func(repo courierDomain.Repository) []uuid.UUID
		expectedError error
	}{
		{
			name: "Success: One courier",
			setup: func(repo courierDomain.Repository) []uuid.UUID {
				return []uuid.UUID{s.createTestCourierInDb(s.createRandomPhone(), repo).ID}
			},
			expectedError: nil,
		},
		{
			name: "Success: Multiple couriers",
			setup: func(repo courierDomain.Repository) []uuid.UUID {
				var courierIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					courier := s.createTestCourierInDb(s.createRandomPhone(), repo)
					courierIDs = append(courierIDs, courier.ID)
				}
				return courierIDs
			},
			expectedError: nil,
		},
		{
			name: "Success: No couriers",
			setup: func(repo courierDomain.Repository) []uuid.UUID {
				return []uuid.UUID{}
			},
			expectedError: nil,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			courierIDs := test.setup(repo)

			couriers, err := repo.GetAll(s.ctx)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				returnedCourierIDs := make([]uuid.UUID, 0, len(couriers))
				for _, courier := range couriers {
					returnedCourierIDs = append(returnedCourierIDs, courier.ID)
				}
				for _, courierID := range courierIDs {
					require.Contains(s.T(), returnedCourierIDs, courierID)
				}
			}
		})
	}
}

func TestCourierRepository(t *testing.T) {
	suite.Run(t, new(CourierRepositoryTestSuite))
}
