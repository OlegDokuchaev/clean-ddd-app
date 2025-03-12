//go:build integration

package publisher

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	domain "warehouse/internal/domain/common"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	outboxPublisher "warehouse/internal/infrastructure/publisher/outbox"
	"warehouse/internal/tests/testutils"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	ProductTopic = "product-topic"
)

type OutboxTestSuite struct {
	suite.Suite
	ctx context.Context

	testMessaging *testutils.TestMessaging

	productWriter *kafka.Writer
	productReader *kafka.Reader
}

func (s *OutboxTestSuite) SetupSuite() {
	s.ctx = context.Background()

	testMessaging, err := testutils.NewTestMessaging(s.ctx)
	require.NoError(s.T(), err)
	s.testMessaging = testMessaging

	err = s.testMessaging.CreateTopics(s.ctx, ProductTopic)
	require.NoError(s.T(), err)

	s.productWriter = s.testMessaging.CreateWriter(ProductTopic)
	s.productReader = s.testMessaging.CreateReader(ProductTopic)
}

func (s *OutboxTestSuite) TearDownSuite() {
	if s.productWriter != nil {
		err := s.productWriter.Close()
		require.NoError(s.T(), err)
		err = s.productReader.Close()
		require.NoError(s.T(), err)
	}

	if s.testMessaging != nil {
		err := s.testMessaging.Container.Terminate(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *OutboxTestSuite) createTestPublisher() outboxDomain.Publisher {
	return outboxPublisher.NewPublisher(s.productWriter)
}

func (s *OutboxTestSuite) TestPublish() {
	tests := []struct {
		name            string
		message         *outboxDomain.Message
		expectedError   error
		validateMessage func(message *outboxDomain.Message, kafkaMsg kafka.Message)
	}{
		{
			name: "Success",
			message: func() *outboxDomain.Message {
				payload := productDomain.CreatedPayload{
					ProductID: uuid.New(),
				}
				event := domain.NewEvent[productDomain.CreatedPayload, productDomain.CreateEvent](payload)
				msg, err := outboxDomain.Create(event)
				require.NoError(s.T(), err)
				return msg
			}(),
			expectedError: nil,
			validateMessage: func(message *outboxDomain.Message, kafkaMsg kafka.Message) {
				var payload outboxDomain.Message
				err := json.Unmarshal(kafkaMsg.Value, &payload)
				require.NoError(s.T(), err)
				require.EqualValues(s.T(), *message, payload)
			},
		},
		{
			name: "Failure: Invalid message type",
			message: &outboxDomain.Message{
				ID:      uuid.New(),
				Type:    "unknown.event",
				Payload: []byte(`{"test": "data"}`),
			},
			validateMessage: func(message *outboxDomain.Message, kafkaMsg kafka.Message) {},
			expectedError:   outboxPublisher.ErrInvalidOutboxMessage,
		},
	}

	publisher := s.createTestPublisher()
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

				kafkaMsg, err := s.productReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				tt.validateMessage(tt.message, kafkaMsg)
			}
		})
	}
}

func TestOutboxTestSuite(t *testing.T) {
	suite.Run(t, new(OutboxTestSuite))
}
