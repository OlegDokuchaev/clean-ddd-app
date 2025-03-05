package domain

import (
	"github.com/google/uuid"
	"testing"
	domain "warehouse/internal/domain/common"
	"warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OutboxDomainTestSuite struct {
	suite.Suite
}

func (suite *OutboxDomainTestSuite) TestCreateOutbox() {
	tests := []struct {
		name        string
		event       domain.Event
		expectedErr error
	}{
		{
			name: "Success",
			event: domain.NewEvent[productDomain.CreatedPayload, productDomain.CreateEvent](
				productDomain.CreatedPayload{
					ProductID: uuid.New(),
				},
			),
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {
			suite.T().Parallel()

			message, err := outbox.Create(tc.event)

			if tc.expectedErr != nil {
				require.Error(suite.T(), err)
				require.ErrorIs(suite.T(), err, tc.expectedErr)
			} else {
				require.NoError(suite.T(), err)
				require.NotNil(suite.T(), message)
			}
		})
	}
}

func TestOutboxDomainTestSuite(t *testing.T) {
	suite.Run(t, new(OutboxDomainTestSuite))
}
