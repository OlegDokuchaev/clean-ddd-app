//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"net"
	createOrder "order/internal/application/order/saga/create_order"
	"order/internal/infrastructure/db"
	"order/internal/infrastructure/db/migrations"
	"order/internal/infrastructure/messaging"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"
	"order/internal/tests/testutils"
	"testing"
	"time"

	appDI "order/internal/application/di"
	infraDI "order/internal/infrastructure/di"
	"order/internal/infrastructure/logger"
	presentationDI "order/internal/presentation/di"
	orderv1 "order/internal/presentation/grpc"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	TestOrderCollectionName = "orders"
	WarehouseTopic          = "warehouse-topic"
	WarehouseResTopic       = "warehouse-topic-res"
	OrderTopic              = "order-topic"
	OrderResTopic           = "order-topic-res"
	CourierTopic            = "courier-topic"
	CourierResTopic         = "courier-topic-res"
)

type CreateOrderE2ESuite struct {
	suite.Suite

	ctx context.Context

	app     *fx.App
	grpcURL string

	db *testutils.TestDB

	testMessaging *testutils.TestMessaging
}

func (s *CreateOrderE2ESuite) SetupSuite() {
	config, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	// 1) MongoDB
	s.db, err = testutils.NewTestDB(s.ctx, config)
	require.NoError(s.T(), err)

	// 2) Kafka
	testMessaging, err := testutils.NewTestMessaging(s.ctx)
	require.NoError(s.T(), err)
	s.testMessaging = testMessaging

	err = s.testMessaging.CreateTopics(s.ctx, WarehouseTopic, OrderTopic, CourierTopic)
	require.NoError(s.T(), err)

	// 3) GRPC
	grpcLn, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err)

	_, grpcPortStr, err := net.SplitHostPort(grpcLn.Addr().String())
	require.NoError(s.T(), err)

	_ = grpcLn.Close()
	s.grpcURL = net.JoinHostPort("127.0.0.1", grpcPortStr)

	// 4) Messaging/DB/GRPC configs
	msgCfg := &messaging.Config{
		Address: s.testMessaging.Address(),

		OrderCmdTopic:           OrderTopic,
		OrderCmdResTopic:        OrderResTopic,
		OrderCmdConsumerGroupID: uuid.NewString(),

		WarehouseCmdTopic:              WarehouseTopic,
		WarehouseCmdResTopic:           WarehouseResTopic,
		WarehouseCmdResConsumerGroupID: uuid.NewString(),

		CourierCmdTopic:              CourierTopic,
		CourierCmdResTopic:           CourierResTopic,
		CourierCmdResConsumerGroupID: uuid.NewString(),
	}

	dbCfg := &db.Config{
		URI:             "", // client & db will be replaced below
		Database:        "name",
		OrderCollection: TestOrderCollectionName,
	}

	grpcCfg := &orderv1.Config{Port: grpcPortStr}

	// 5) Logrus
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	// 5) FX app
	s.app = fx.New(
		infraDI.LoggerModule,
		infraDI.MessagingModule,
		infraDI.DatabaseModule,
		infraDI.RepositoryModule,
		infraDI.PublisherModule,
		appDI.UseCaseModule,
		appDI.SagaModule,
		presentationDI.GRPCModule,
		fx.Replace(log),
		fx.Replace(msgCfg),
		fx.Replace(dbCfg),
		fx.Replace(grpcCfg),
		fx.Replace(s.db.DB.Client()),
		fx.Replace(s.db.DB),
		fx.Invoke(func(lc fx.Lifecycle, l logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error { return nil },
				OnStop:  func(context.Context) error { return nil },
			})
		}),
	)

	startCtx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()
	require.NoError(s.T(), s.app.Start(startCtx))

	// 6) Wait until gRPC is ready (TCP connect)
	require.Eventually(s.T(), func() bool {
		c, err := net.DialTimeout("tcp", s.grpcURL, 2*time.Second)
		if err != nil {
			return false
		}
		_ = c.Close()
		return true
	}, 10*time.Second, 200*time.Millisecond)
}

func (s *CreateOrderE2ESuite) TearDownSuite() {
	if s.app != nil {
		ctx, cancel := context.WithTimeout(s.ctx, 20*time.Second)
		_ = s.app.Stop(ctx)
		cancel()
	}

	if s.db != nil {
		err := s.db.Close(s.ctx)
		require.NoError(s.T(), err)
	}
	if s.testMessaging != nil {
		err := s.testMessaging.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *CreateOrderE2ESuite) Test_CreateOrder_Success() {
	// 1) Dial gRPC client
	conn, err := grpc.NewClient(s.grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.T(), err)
	defer func() { _ = conn.Close() }()

	client := orderv1.NewOrderServiceClient(conn)

	// 2) Build request
	customerID := uuid.New().String()
	productID := uuid.New().String()
	req := &orderv1.CreateOrderRequest{
		CustomerId: customerID,
		Address:    "Some Address",
		Items: []*orderv1.OrderItem{
			{ProductId: productID, Price: 100.0, Count: 1},
		},
	}

	// 3) Call CreateOrder
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()
	res, err := client.CreateOrder(ctx, req)

	// 5) Assert result
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.GetOrderId())

	// 5) Assert DB contains the order
	orders := s.db.DB.Collection(TestOrderCollectionName)
	var doc bson.M
	findCtx, findCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer findCancel()
	err = orders.FindOne(findCtx, bson.M{"_id": res.GetOrderId()}).Decode(&doc)
	require.NoError(s.T(), err)

	// 6) Assert Kafka saga message (ReserveItemsCmd) published to warehouse-topic
	reader := s.testMessaging.CreateReader(WarehouseTopic)
	defer func() { _ = reader.Close() }()

	readCtx, readCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer readCancel()
	msg, err := reader.ReadMessage(readCtx)
	require.NoError(s.T(), err)

	var cmdMessage createOrderPublisher.CmdMessage
	require.NoError(s.T(), json.Unmarshal(msg.Value, &cmdMessage))
	require.Equal(s.T(), createOrderPublisher.ReserveItemsCmdName, cmdMessage.Name)

	payloadRaw, err := json.Marshal(cmdMessage.Payload)
	require.NoError(s.T(), err)
	var payload createOrder.ReserveItemsCmd
	require.NoError(s.T(), json.Unmarshal(payloadRaw, &payload))

	createdID := uuid.MustParse(res.GetOrderId())
	require.Equal(s.T(), createdID, payload.OrderID)
}

func (s *CreateOrderE2ESuite) Test_CreateOrder_InvalidData() {
	// 1) Dial gRPC client
	conn, err := grpc.NewClient(s.grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.T(), err)
	defer func() { _ = conn.Close() }()

	client := orderv1.NewOrderServiceClient(conn)

	// 2) Build invalid request (empty address)
	customerID := uuid.New().String()
	productID := uuid.New().String()
	req := &orderv1.CreateOrderRequest{
		CustomerId: customerID,
		Address:    "", // invalid
		Items: []*orderv1.OrderItem{
			{ProductId: productID, Price: 100.0, Count: 1},
		},
	}

	// 3) Call CreateOrder and expect InvalidArgument
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	res, err := client.CreateOrder(ctx, req)

	// 4) Assert result
	require.Error(s.T(), err)
	_ = res
	st, _ := status.FromError(err)
	require.Equal(s.T(), codes.InvalidArgument, st.Code())
}

func Test_CreateOrderE2E(t *testing.T) {
	suite.Run(t, new(CreateOrderE2ESuite))
}
