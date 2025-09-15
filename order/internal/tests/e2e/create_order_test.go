//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"net"
	createOrder "order/internal/application/order/saga/create_order"
	"order/internal/infrastructure/db/migrations"
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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type CreateOrderE2ESuite struct {
	suite.Suite

	ctx context.Context

	app     *fx.App
	grpcURL string

	db        *testutils.TestDB
	messaging *testutils.TestMessaging
}

func (s *CreateOrderE2ESuite) BeforeAll(t provider.T) {
	config, err := migrations.NewConfig()
	t.Require().NoError(err)

	s.ctx = context.Background()

	// 1) MongoDB
	s.db, err = testutils.NewTestDB(s.ctx, config, nil)
	t.Require().NoError(err)

	// 2) Kafka
	testMessaging, err := testutils.NewTestMessaging(s.ctx, nil)
	t.Require().NoError(err)
	s.messaging = testMessaging

	// 3) GRPC
	grpcLn, err := net.Listen("tcp", "127.0.0.1:0")
	t.Require().NoError(err)

	_, grpcPortStr, err := net.SplitHostPort(grpcLn.Addr().String())
	t.Require().NoError(err)

	_ = grpcLn.Close()
	s.grpcURL = net.JoinHostPort("127.0.0.1", grpcPortStr)

	grpcCfg := &orderv1.Config{Port: grpcPortStr}

	// 4) Logrus
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
		fx.Replace(logrus.New()),
		fx.Replace(s.messaging.Cfg),
		fx.Replace(s.db.Cfg),
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

	err = s.app.Start(startCtx)
	t.Require().NoError(err)

	// 6) Wait until gRPC is ready (TCP connect)
	t.Require().Eventually(func() bool {
		c, err := net.DialTimeout("tcp", s.grpcURL, 2*time.Second)
		if err != nil {
			return false
		}
		_ = c.Close()
		return true
	}, 10*time.Second, 200*time.Millisecond)
}

func (s *CreateOrderE2ESuite) AfterAll(t provider.T) {
	if s.app != nil {
		ctx, cancel := context.WithTimeout(s.ctx, 20*time.Second)
		_ = s.app.Stop(ctx)
		cancel()
	}

	if s.db != nil {
		err := s.db.Close(s.ctx)
		t.Require().NoError(err)
	}
	if s.messaging != nil {
		err := s.messaging.Close(s.ctx)
		t.Require().NoError(err)
	}
}

func (s *CreateOrderE2ESuite) AfterEach(t provider.T) {
	err := s.db.Clear(s.ctx)
	t.Require().NoError(err)

	err = s.messaging.Clear(s.ctx)
	t.Require().NoError(err)
}

func (s *CreateOrderE2ESuite) Test_CreateOrder_Success(t provider.T) {
	// 1) Dial gRPC client
	conn, err := grpc.NewClient(s.grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Require().NoError(err)
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
	t.Require().NoError(err)
	t.Require().NotEmpty(res.GetOrderId())

	// 5) Assert DB contains the order
	orders := s.db.DB.Collection(s.db.Cfg.OrderCollection)
	var doc bson.M
	findCtx, findCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer findCancel()
	err = orders.FindOne(findCtx, bson.M{"_id": res.GetOrderId()}).Decode(&doc)
	t.Require().NoError(err)

	// 6) Assert Kafka saga message (ReserveItemsCmd) published to warehouse-topic
	reader := s.messaging.CreateReader(s.messaging.Cfg.WarehouseCmdTopic)
	defer func() { _ = reader.Close() }()

	readCtx, readCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer readCancel()
	msg, err := reader.ReadMessage(readCtx)
	t.Require().NoError(err)

	var cmdMessage createOrderPublisher.CmdMessage
	t.Require().NoError(json.Unmarshal(msg.Value, &cmdMessage))
	t.Require().Equal(createOrderPublisher.ReserveItemsCmdName, cmdMessage.Name)

	payloadRaw, err := json.Marshal(cmdMessage.Payload)
	t.Require().NoError(err)
	var payload createOrder.ReserveItemsCmd
	t.Require().NoError(json.Unmarshal(payloadRaw, &payload))

	createdID := uuid.MustParse(res.GetOrderId())
	t.Require().Equal(createdID, payload.OrderID)
}

func (s *CreateOrderE2ESuite) Test_CreateOrder_InvalidData(t provider.T) {
	// 1) Dial gRPC client
	conn, err := grpc.NewClient(s.grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Require().NoError(err)
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
	t.Require().Error(err)
	_ = res
	st, _ := status.FromError(err)
	t.Require().Equal(codes.InvalidArgument, st.Code())
}

func Test_CreateOrderE2E(t *testing.T) {
	suite.RunSuite(t, new(CreateOrderE2ESuite))
}
