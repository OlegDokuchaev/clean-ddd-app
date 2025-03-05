package courier

import (
	"context"
	courierApplication "courier/internal/application/courier"
	courierDomain "courier/internal/domain/courier"
	courierMock "courier/internal/mocks/courier"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CourierUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *CourierUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *CourierUseCaseTestSuite) createTestCourier() *courierDomain.Courier {
	courier, err := courierDomain.Create("test", "+79032895555", "password")
	require.NoError(s.T(), err)
	return courier
}

func (s *CourierUseCaseTestSuite) createTestCouriers(count int) []*courierDomain.Courier {
	couriers := make([]*courierDomain.Courier, 0, count)
	for range count {
		couriers = append(couriers, s.createTestCourier())
	}
	return couriers
}

func (s *CourierUseCaseTestSuite) TestAssignOrder() {
	tests := []struct {
		name        string
		orderID     uuid.UUID
		setup       func(repo *courierMock.RepositoryMock) []uuid.UUID
		expectedErr error
	}{
		{
			name:    "Success: Multiple couriers found",
			orderID: uuid.New(),
			setup: func(repo *courierMock.RepositoryMock) []uuid.UUID {
				couriers := s.createTestCouriers(10)
				repo.On("GetAll", s.ctx).Return(couriers, nil).Once()

				courierIDs := make([]uuid.UUID, 0, len(couriers))
				for _, courier := range couriers {
					courierIDs = append(courierIDs, courier.ID)
				}
				return courierIDs
			},
			expectedErr: nil,
		},
		{
			name:    "Success: One courier found",
			orderID: uuid.New(),
			setup: func(repo *courierMock.RepositoryMock) []uuid.UUID {
				courier := s.createTestCourier()
				repo.On("GetAll", s.ctx).Return([]*courierDomain.Courier{courier}, nil).Once()
				return []uuid.UUID{courier.ID}
			},
			expectedErr: nil,
		},
		{
			name:    "Failure: Available courier not found",
			orderID: uuid.New(),
			setup: func(repo *courierMock.RepositoryMock) []uuid.UUID {
				repo.On("GetAll", s.ctx).Return([]*courierDomain.Courier{}, nil).Once()
				return []uuid.UUID{}
			},
			expectedErr: courierApplication.ErrAvailableCourierNotFound,
		},
		{
			name:    "Failure: Courier repository get all error",
			orderID: uuid.New(),
			setup: func(repo *courierMock.RepositoryMock) []uuid.UUID {
				repo.On("GetAll", s.ctx).
					Return([]*courierDomain.Courier{}, errors.New("get all couriers error")).Once()
				return []uuid.UUID{}
			},
			expectedErr: errors.New("get all couriers error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(courierMock.RepositoryMock)
			uc := courierApplication.NewUseCase(repo)
			courierIDs := tc.setup(repo)

			courierID, err := uc.AssignOrder(s.ctx, tc.orderID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Contains(s.T(), courierIDs, courierID)

			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Equal(s.T(), uuid.Nil, courierID)
			}

			repo.AssertExpectations(s.T())
		})
	}
}

func TestCourierUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CourierUseCaseTestSuite))
}
