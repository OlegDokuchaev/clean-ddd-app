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
		expectedError   error
		validateMessage func(message *outbox.Message, kafkaMsg kafka.Message)
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
			expectedError: nil,
			validateMessage: func(message *outbox.Message, kafkaMsg kafka.Message) {
				var payload outbox.Message
				err := json.Unmarshal(kafkaMsg.Value, &payload)
				require.NoError(s.T(), err)
				require.EqualValues(s.T(), *message, payload)
			},
		},
		{
			name: "Failure: Invalid message type",
			message: &outbox.Message{
				ID:      uuid.New(),
				Type:    "unknown.event",
				Payload: []byte(`{"test": "data"}`),
			},
			validateMessage: func(message *outbox.Message, kafkaMsg kafka.Message) {},
			expectedError:   outboxPublisher.ErrInvalidOutboxMessage,
		},
	}

	publisher := outboxPublisher.NewPublisher(s.testMessaging.Writer)
	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := publisher.Publish(s.ctx, tt.message)

			if tt.expectedError != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.expectedError)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				kafkaMsg, err := s.testMessaging.Reader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				tt.validateMessage(tt.message, kafkaMsg)
			}
		})
	}
}

func TestOutboxTestSuite(t *testing.T) {
	suite.Run(t, new(OutboxTestSuite))
}
