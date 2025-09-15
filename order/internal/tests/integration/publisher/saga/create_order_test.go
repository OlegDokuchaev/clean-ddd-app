//go:build integration

package saga

import (
	"context"
	"encoding/json"
	createOrder "order/internal/application/order/saga/create_order"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"
	"order/internal/tests/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/segmentio/kafka-go"
)

type CreateOrderPublisherTestSuite struct {
	suite.Suite
	ctx context.Context

	messaging *testutils.TestMessaging

	warehouseWriter *kafka.Writer
	warehouseReader *kafka.Reader

	orderWriter *kafka.Writer
	orderReader *kafka.Reader

	courierWriter *kafka.Writer
	courierReader *kafka.Reader
}

func (s *CreateOrderPublisherTestSuite) BeforeAll(t provider.T) {
	s.ctx = context.Background()

	testMessaging, err := testutils.NewTestMessaging(s.ctx, nil)
	t.Require().NoError(err)
	s.messaging = testMessaging
}

func (s *CreateOrderPublisherTestSuite) AfterAll(t provider.T) {
	if s.messaging != nil {
		err := s.messaging.Close(s.ctx)
		t.Require().NoError(err)
	}
}

func (s *CreateOrderPublisherTestSuite) BeforeEach(_ provider.T) {
	s.warehouseWriter = s.messaging.CreateWriter(s.messaging.Cfg.WarehouseCmdTopic)
	s.warehouseReader = s.messaging.CreateReader(s.messaging.Cfg.WarehouseCmdTopic)

	s.orderWriter = s.messaging.CreateWriter(s.messaging.Cfg.OrderCmdTopic)
	s.orderReader = s.messaging.CreateReader(s.messaging.Cfg.OrderCmdTopic)

	s.courierWriter = s.messaging.CreateWriter(s.messaging.Cfg.CourierCmdTopic)
	s.courierReader = s.messaging.CreateReader(s.messaging.Cfg.CourierCmdTopic)
}

func (s *CreateOrderPublisherTestSuite) AfterEach(t provider.T) {
	if s.warehouseWriter != nil {
		err := s.warehouseWriter.Close()
		t.Require().NoError(err)
		err = s.warehouseReader.Close()
		t.Require().NoError(err)
	}

	if s.orderWriter != nil {
		err := s.orderWriter.Close()
		t.Require().NoError(err)
		err = s.orderReader.Close()
		t.Require().NoError(err)
	}

	if s.courierWriter != nil {
		err := s.courierWriter.Close()
		t.Require().NoError(err)
		err = s.courierReader.Close()
		t.Require().NoError(err)
	}

	err := s.messaging.Clear(s.ctx)
	t.Require().NoError(err)
}

func (s *CreateOrderPublisherTestSuite) createTestPublisher() createOrder.Publisher {
	return createOrderPublisher.NewPublisher(s.warehouseWriter, s.orderWriter, s.courierWriter)
}

func (s *CreateOrderPublisherTestSuite) TestPublishReserveItemsCmd(t provider.T) {
	tests := []struct {
		name            string
		cmd             createOrder.ReserveItemsCmd
		validateMessage func(cmd createOrder.ReserveItemsCmd, message kafka.Message)
		expectedError   error
	}{
		{
			name: "Success",
			cmd: createOrder.ReserveItemsCmd{
				OrderID: uuid.New(),
				Items: []createOrder.OrderItem{
					{
						ProductID: uuid.New(),
						Count:     1,
					},
				},
			},
			validateMessage: func(cmd createOrder.ReserveItemsCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				t.Require().NoError(err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				t.Require().NoError(err)

				var payload createOrder.ReserveItemsCmd
				err = json.Unmarshal(encodedPayload, &payload)
				t.Require().NoError(err)

				t.Require().EqualValues(cmd, payload)
			},
			expectedError: nil,
		},
	}

	publisher := s.createTestPublisher()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t provider.T) {
			err := publisher.PublishReserveItemsCmd(s.ctx, tt.cmd)

			if tt.expectedError != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tt.expectedError)
			} else {
				t.Require().NoError(err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.warehouseReader.ReadMessage(ctx)
				t.Require().NoError(err)
				tt.validateMessage(tt.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishReleaseItemsCmd(t provider.T) {
	tests := []struct {
		name            string
		cmd             createOrder.ReleaseItemsCmd
		validateMessage func(cmd createOrder.ReleaseItemsCmd, message kafka.Message)
		expectedError   error
	}{
		{
			name: "Success",
			cmd: createOrder.ReleaseItemsCmd{
				OrderID: uuid.New(),
				Items: []createOrder.OrderItem{
					{
						ProductID: uuid.New(),
						Count:     1,
					},
				},
			},
			validateMessage: func(cmd createOrder.ReleaseItemsCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				t.Require().NoError(err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				t.Require().NoError(err)

				var payload createOrder.ReleaseItemsCmd
				err = json.Unmarshal(encodedPayload, &payload)
				t.Require().NoError(err)

				t.Require().EqualValues(cmd, payload)
			},
			expectedError: nil,
		},
	}

	publisher := s.createTestPublisher()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t provider.T) {
			err := publisher.PublishReleaseItemsCmd(s.ctx, tt.cmd)

			if tt.expectedError != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tt.expectedError)
			} else {
				t.Require().NoError(err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.warehouseReader.ReadMessage(ctx)
				t.Require().NoError(err)
				tt.validateMessage(tt.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishCancelOutOfStockCmd(t provider.T) {
	tests := []struct {
		name            string
		cmd             createOrder.CancelOutOfStockCmd
		validateMessage func(cmd createOrder.CancelOutOfStockCmd, message kafka.Message)
		expectedError   error
	}{
		{
			name: "Success",
			cmd: createOrder.CancelOutOfStockCmd{
				OrderID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.CancelOutOfStockCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				t.Require().NoError(err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				t.Require().NoError(err)

				var payload createOrder.CancelOutOfStockCmd
				err = json.Unmarshal(encodedPayload, &payload)
				t.Require().NoError(err)

				t.Require().EqualValues(cmd, payload)
			},
			expectedError: nil,
		},
	}

	publisher := s.createTestPublisher()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t provider.T) {
			err := publisher.PublishCancelOutOfStockCmd(s.ctx, tt.cmd)

			if tt.expectedError != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tt.expectedError)
			} else {
				t.Require().NoError(err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.orderReader.ReadMessage(ctx)
				t.Require().NoError(err)
				tt.validateMessage(tt.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishAssignCourierCmd(t provider.T) {
	tests := []struct {
		name            string
		cmd             createOrder.AssignCourierCmd
		validateMessage func(cmd createOrder.AssignCourierCmd, message kafka.Message)
		expectedError   error
	}{
		{
			name: "Success",
			cmd: createOrder.AssignCourierCmd{
				OrderID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.AssignCourierCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				t.Require().NoError(err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				t.Require().NoError(err)

				var payload createOrder.AssignCourierCmd
				err = json.Unmarshal(encodedPayload, &payload)
				t.Require().NoError(err)

				t.Require().EqualValues(cmd, payload)
			},
			expectedError: nil,
		},
	}

	publisher := s.createTestPublisher()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t provider.T) {
			err := publisher.PublishAssignCourierCmd(s.ctx, tt.cmd)

			if tt.expectedError != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tt.expectedError)
			} else {
				t.Require().NoError(err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.courierReader.ReadMessage(ctx)
				t.Require().NoError(err)
				tt.validateMessage(tt.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishBeginDeliveryCmd(t provider.T) {
	tests := []struct {
		name            string
		cmd             createOrder.BeginDeliveryCmd
		validateMessage func(cmd createOrder.BeginDeliveryCmd, message kafka.Message)
		expectedError   error
	}{
		{
			name: "Success",
			cmd: createOrder.BeginDeliveryCmd{
				OrderID:   uuid.New(),
				CourierID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.BeginDeliveryCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				t.Require().NoError(err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				t.Require().NoError(err)

				var payload createOrder.BeginDeliveryCmd
				err = json.Unmarshal(encodedPayload, &payload)
				t.Require().NoError(err)

				t.Require().EqualValues(cmd, payload)
			},
			expectedError: nil,
		},
	}

	publisher := s.createTestPublisher()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t provider.T) {
			err := publisher.PublishBeginDeliveryCmd(s.ctx, tt.cmd)

			if tt.expectedError != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tt.expectedError)
			} else {
				t.Require().NoError(err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.orderReader.ReadMessage(ctx)
				t.Require().NoError(err)
				tt.validateMessage(tt.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishCancelCourierNotFoundCmd(t provider.T) {
	tests := []struct {
		name            string
		cmd             createOrder.CancelCourierNotFoundCmd
		validateMessage func(cmd createOrder.CancelCourierNotFoundCmd, message kafka.Message)
		expectedError   error
	}{
		{
			name: "Success",
			cmd: createOrder.CancelCourierNotFoundCmd{
				OrderID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.CancelCourierNotFoundCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				t.Require().NoError(err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				t.Require().NoError(err)

				var payload createOrder.CancelCourierNotFoundCmd
				err = json.Unmarshal(encodedPayload, &payload)
				t.Require().NoError(err)

				t.Require().EqualValues(cmd, payload)
			},
			expectedError: nil,
		},
	}

	publisher := s.createTestPublisher()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t provider.T) {
			err := publisher.PublishCancelCourierNotFoundCmd(s.ctx, tt.cmd)

			if tt.expectedError != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tt.expectedError)
			} else {
				t.Require().NoError(err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.orderReader.ReadMessage(ctx)
				t.Require().NoError(err)
				tt.validateMessage(tt.cmd, message)
			}
		})
	}
}

func TestCreateOrderPublisherTestSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateOrderPublisherTestSuite))
}
