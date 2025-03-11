//go:build integration

package publisher

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	domain "warehouse/internal/domain/common"
	"warehouse/internal/domain/outbox"
	"warehouse/internal/domain/product"
	outboxPublisher "warehouse/internal/infrastructure/publisher/outbox"
	"warehouse/internal/tests/testutils"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OutboxTestSuite struct {
	suite.Suite
	ctx           context.Context
	testMessaging *testutils.TestMessaging
}

func (s *OutboxTestSuite) SetupSuite() {
	s.ctx = context.Background()

	testMessaging, err := testutils.NewTestMessaging(s.ctx)
	require.NoError(s.T(), err)
	s.testMessaging = testMessaging
}

func (s *OutboxTestSuite) TearDownSuite() {
	if s.testMessaging != nil {
		err := s.testMessaging.Writer.Close()
		require.NoError(s.T(), err)

		err = s.testMessaging.Reader.Close()
		require.NoError(s.T(), err)

		err = s.testMessaging.Container.Terminate(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *OutboxTestSuite) TestPublish() {
	tests := []struct {
		name            string
		message         *outbox.Message
		setupPublisher  func() outbox.Publisher
		expectedError   error
		validateMessage func(t require.TestingT, message kafka.Message)
	}{
		{
			name: "Success",
			message: func() *outbox.Message {
				event := domain.NewEvent[product.CreatedPayload, product.CreateEvent](product.CreatedPayload{
					ProductID: uuid.New(),
				})
				msg, err := outbox.Create(event)
				require.NoError(s.T(), err)
				return msg
			}(),
			setupPublisher: func() outbox.Publisher {
				config := &outboxPublisher.Config{ProductTopic: testutils.TestTopic}
				return outboxPublisher.NewPublisher(config, s.testMessaging.Writer)
			},
			expectedError: nil,
			validateMessage: func(t require.TestingT, message kafka.Message) {
				var payload product.CreatedPayload
				err := json.Unmarshal(message.Value, &payload)
				require.NoError(t, err)
				require.NotEmpty(t, payload.ProductID)
			},
		},
		{
			name: "Failure: Invalid message type",
			message: &outbox.Message{
				ID:      uuid.New(),
				Type:    "unknown.event",
				Payload: []byte(`{"test": "data"}`),
			},
			setupPublisher: func() outbox.Publisher {
				config := &outboxPublisher.Config{ProductTopic: testutils.TestTopic}
				return outboxPublisher.NewPublisher(config, s.testMessaging.Writer)
			},
			expectedError: outboxPublisher.ErrInvalidOutboxMessage,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			publisher := tt.setupPublisher()

			err := publisher.Publish(s.ctx, tt.message)

			if tt.expectedError != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.expectedError)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.testMessaging.Reader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				tt.validateMessage(s.T(), message)
			}
		})
	}
}

func TestOutboxTestSuite(t *testing.T) {
	suite.Run(t, new(OutboxTestSuite))
}
