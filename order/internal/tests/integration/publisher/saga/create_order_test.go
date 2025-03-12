//go:build integration

package saga

import (
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"
	"order/internal/tests/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	WarehouseTopic = "warehouse-topic"
	OrderTopic     = "order-topic"
	CourierTopic   = "courier-topic"
)

type CreateOrderPublisherTestSuite struct {
	suite.Suite
	ctx context.Context

	testMessaging *testutils.TestMessaging

	warehouseWriter *kafka.Writer
	warehouseReader *kafka.Reader

	orderWriter *kafka.Writer
	orderReader *kafka.Reader

	courierWriter *kafka.Writer
	courierReader *kafka.Reader
}

func (s *CreateOrderPublisherTestSuite) SetupSuite() {
	s.ctx = context.Background()

	testMessaging, err := testutils.NewTestMessaging(s.ctx)
	require.NoError(s.T(), err)
	s.testMessaging = testMessaging

	err = s.testMessaging.CreateTopics(s.ctx, WarehouseTopic, OrderTopic, CourierTopic)
	require.NoError(s.T(), err)

	s.warehouseWriter = s.testMessaging.CreateWriter(WarehouseTopic)
	s.warehouseReader = s.testMessaging.CreateReader(WarehouseTopic)

	s.orderWriter = s.testMessaging.CreateWriter(OrderTopic)
	s.orderReader = s.testMessaging.CreateReader(OrderTopic)

	s.courierWriter = s.testMessaging.CreateWriter(CourierTopic)
	s.courierReader = s.testMessaging.CreateReader(CourierTopic)
}

func (s *CreateOrderPublisherTestSuite) TearDownSuite() {
	if s.warehouseWriter != nil {
		err := s.warehouseWriter.Close()
		require.NoError(s.T(), err)
		err = s.warehouseReader.Close()
		require.NoError(s.T(), err)
	}

	if s.orderWriter != nil {
		err := s.orderWriter.Close()
		require.NoError(s.T(), err)
		err = s.orderReader.Close()
		require.NoError(s.T(), err)
	}

	if s.courierWriter != nil {
		err := s.courierWriter.Close()
		require.NoError(s.T(), err)
		err = s.courierReader.Close()
		require.NoError(s.T(), err)
	}

	if s.testMessaging != nil {
		err := s.testMessaging.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *CreateOrderPublisherTestSuite) createTestPublisher() createOrder.Publisher {
	return createOrderPublisher.NewPublisher(s.warehouseWriter, s.orderWriter, s.courierWriter)
}

func (s *CreateOrderPublisherTestSuite) TestPublishReserveItemsCmd() {
	tests := []struct {
		name            string
		cmd             createOrder.ReserveItemsCmd
		validateMessage func(cmd createOrder.ReserveItemsCmd, message kafka.Message)
		expectedError   bool
	}{
		{
			name: "Success",
			cmd: createOrder.ReserveItemsCmd{
				OrderID: uuid.New(),
				Items: []orderDomain.Item{
					{
						ProductID: uuid.New(),
						Price:     decimal.NewFromInt(100),
						Count:     1,
					},
				},
			},
			validateMessage: func(cmd createOrder.ReserveItemsCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				require.NoError(s.T(), err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				require.NoError(s.T(), err)

				var payload createOrder.ReserveItemsCmd
				err = json.Unmarshal(encodedPayload, &payload)
				require.NoError(s.T(), err)

				require.EqualValues(s.T(), cmd, payload)
			},
			expectedError: false,
		},
	}

	publisher := s.createTestPublisher()
	for _, test := range tests {
		s.Run(test.name, func() {
			err := publisher.PublishReserveItemsCmd(s.ctx, test.cmd)

			if test.expectedError {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.warehouseReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				test.validateMessage(test.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishReleaseItemsCmd() {
	tests := []struct {
		name            string
		cmd             createOrder.ReleaseItemsCmd
		validateMessage func(cmd createOrder.ReleaseItemsCmd, message kafka.Message)
		expectedError   bool
	}{
		{
			name: "Success",
			cmd: createOrder.ReleaseItemsCmd{
				OrderID: uuid.New(),
				Items: []orderDomain.Item{
					{
						ProductID: uuid.New(),
						Price:     decimal.NewFromInt(100),
						Count:     1,
					},
				},
			},
			validateMessage: func(cmd createOrder.ReleaseItemsCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				require.NoError(s.T(), err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				require.NoError(s.T(), err)

				var payload createOrder.ReleaseItemsCmd
				err = json.Unmarshal(encodedPayload, &payload)
				require.NoError(s.T(), err)

				require.EqualValues(s.T(), cmd, payload)
			},
			expectedError: false,
		},
	}

	publisher := s.createTestPublisher()
	for _, test := range tests {
		s.Run(test.name, func() {
			err := publisher.PublishReleaseItemsCmd(s.ctx, test.cmd)

			if test.expectedError {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.warehouseReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				test.validateMessage(test.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishCancelOutOfStockCmd() {
	tests := []struct {
		name            string
		cmd             createOrder.CancelOutOfStockCmd
		validateMessage func(cmd createOrder.CancelOutOfStockCmd, message kafka.Message)
		expectedError   bool
	}{
		{
			name: "Success",
			cmd: createOrder.CancelOutOfStockCmd{
				OrderID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.CancelOutOfStockCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				require.NoError(s.T(), err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				require.NoError(s.T(), err)

				var payload createOrder.CancelOutOfStockCmd
				err = json.Unmarshal(encodedPayload, &payload)
				require.NoError(s.T(), err)

				require.EqualValues(s.T(), cmd, payload)
			},
			expectedError: false,
		},
	}

	publisher := s.createTestPublisher()
	for _, test := range tests {
		s.Run(test.name, func() {
			err := publisher.PublishCancelOutOfStockCmd(s.ctx, test.cmd)

			if test.expectedError {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.orderReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				test.validateMessage(test.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishAssignCourierCmd() {
	tests := []struct {
		name            string
		cmd             createOrder.AssignCourierCmd
		validateMessage func(cmd createOrder.AssignCourierCmd, message kafka.Message)
		expectedError   bool
	}{
		{
			name: "Success",
			cmd: createOrder.AssignCourierCmd{
				OrderID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.AssignCourierCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				require.NoError(s.T(), err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				require.NoError(s.T(), err)

				var payload createOrder.AssignCourierCmd
				err = json.Unmarshal(encodedPayload, &payload)
				require.NoError(s.T(), err)

				require.EqualValues(s.T(), cmd, payload)
			},
			expectedError: false,
		},
	}

	publisher := s.createTestPublisher()
	for _, test := range tests {
		s.Run(test.name, func() {
			err := publisher.PublishAssignCourierCmd(s.ctx, test.cmd)

			if test.expectedError {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.courierReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				test.validateMessage(test.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishBeginDeliveryCmd() {
	tests := []struct {
		name            string
		cmd             createOrder.BeginDeliveryCmd
		validateMessage func(cmd createOrder.BeginDeliveryCmd, message kafka.Message)
		expectedError   bool
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
				require.NoError(s.T(), err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				require.NoError(s.T(), err)

				var payload createOrder.BeginDeliveryCmd
				err = json.Unmarshal(encodedPayload, &payload)
				require.NoError(s.T(), err)

				require.EqualValues(s.T(), cmd, payload)
			},
			expectedError: false,
		},
	}

	publisher := s.createTestPublisher()
	for _, test := range tests {
		s.Run(test.name, func() {
			err := publisher.PublishBeginDeliveryCmd(s.ctx, test.cmd)

			if test.expectedError {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.orderReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				test.validateMessage(test.cmd, message)
			}
		})
	}
}

func (s *CreateOrderPublisherTestSuite) TestPublishCancelCourierNotFoundCmd() {
	tests := []struct {
		name            string
		cmd             createOrder.CancelCourierNotFoundCmd
		validateMessage func(cmd createOrder.CancelCourierNotFoundCmd, message kafka.Message)
		expectedError   bool
	}{
		{
			name: "Success",
			cmd: createOrder.CancelCourierNotFoundCmd{
				OrderID: uuid.New(),
			},
			validateMessage: func(cmd createOrder.CancelCourierNotFoundCmd, message kafka.Message) {
				var cmdMessage createOrderPublisher.CmdMessage
				err := json.Unmarshal(message.Value, &cmdMessage)
				require.NoError(s.T(), err)

				encodedPayload, err := json.Marshal(cmdMessage.Payload)
				require.NoError(s.T(), err)

				var payload createOrder.CancelCourierNotFoundCmd
				err = json.Unmarshal(encodedPayload, &payload)
				require.NoError(s.T(), err)

				require.EqualValues(s.T(), cmd, payload)
			},
			expectedError: false,
		},
	}

	publisher := s.createTestPublisher()
	for _, test := range tests {
		s.Run(test.name, func() {
			err := publisher.PublishCancelCourierNotFoundCmd(s.ctx, test.cmd)

			if test.expectedError {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)

				ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
				defer cancel()

				message, err := s.orderReader.ReadMessage(ctx)
				require.NoError(s.T(), err)
				test.validateMessage(test.cmd, message)
			}
		})
	}
}

func TestCreateOrderPublisherTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderPublisherTestSuite))
}
