//go:build e2e

package e2e

import (
	"context"
	"crypto/tls"
	"customer/internal/tests/testutils/mothers"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	appDI "customer/internal/application/di"
	"customer/internal/infrastructure/db/migrations"
	infraDI "customer/internal/infrastructure/di"
	"customer/internal/infrastructure/logger"
	customerRepository "customer/internal/infrastructure/repository/customer"
	presentationDI "customer/internal/presentation/di"
	customerv1 "customer/internal/presentation/grpc"
	"customer/internal/tests/testutils"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginE2ESuite struct {
	suite.Suite

	ctx context.Context

	app     *fx.App
	tCfg    *testutils.Config
	grpcURL string

	db    *testutils.TestDB
	redis *testutils.TestRedis
}

func (s *LoginE2ESuite) SetupSuite() {
	tCfg, err := testutils.NewConfig()
	require.NoError(s.T(), err)
	s.tCfg = tCfg

	mCfg, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	// 1) Test db
	s.db, err = testutils.NewTestDB(s.ctx, s.tCfg, mCfg)
	require.NoError(s.T(), err)

	// 2) Test redis
	s.redis, err = testutils.NewTestRedis(s.ctx, s.tCfg)
	require.NoError(s.T(), err)

	s.clear()

	// 2) GRPC
	grpcLn, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err)

	_, grpcPortStr, err := net.SplitHostPort(grpcLn.Addr().String())
	require.NoError(s.T(), err)

	_ = grpcLn.Close()
	s.grpcURL = net.JoinHostPort("127.0.0.1", grpcPortStr)

	grpcCfg := &customerv1.Config{Port: grpcPortStr}

	// 4) Logrus
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	// 5) FX app
	s.app = fx.New(
		infraDI.LoggerModule,
		infraDI.DatabaseModule,
		infraDI.RepositoryModule,
		infraDI.TokenManagerModule,
		infraDI.MailSenderModule,
		infraDI.OtpStoreModule,
		infraDI.AuthPoliciesModule,
		appDI.UseCaseModule,
		presentationDI.GRPCModule,
		fx.Replace(logrus.New()),
		fx.Replace(s.db.Cfg),
		fx.Replace(s.db.DB),
		fx.Replace(s.redis.Cfg),
		fx.Replace(s.redis.Client),
		fx.Replace(grpcCfg),
		fx.Invoke(func(lc fx.Lifecycle, l logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error { return nil },
				OnStop:  func(context.Context) error { return nil },
			})
		}),
	)

	startCtx, cancel := context.WithTimeout(s.ctx, time.Minute)
	defer cancel()

	err = s.app.Start(startCtx)
	require.NoError(s.T(), err)

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

func (s *LoginE2ESuite) TearDownSuite() {
	if s.app != nil {
		ctx, cancel := context.WithTimeout(s.ctx, 20*time.Second)
		_ = s.app.Stop(ctx)
		cancel()
	}

	s.clear()
}

func (s *LoginE2ESuite) AfterTest(_, _ string) {
	s.clear()
}

func (s *LoginE2ESuite) clear() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	if s.db != nil {
		err := s.db.Clear(ctx)
		require.NoError(s.T(), err)
	}

	if s.redis != nil {
		err := s.redis.Clear(ctx)
		require.NoError(s.T(), err)
	}
}

func (s *LoginE2ESuite) Test_Login_SendsRealEmailAndReturnsChallenge() {
	// Prepare customer in DB
	customer := mothers.CustomerWithPassword("P@ssw0rd!")
	customer.Email = s.tCfg.ImapUsername

	r := customerRepository.New(s.db.DB)
	require.NoError(s.T(), r.Create(s.ctx, customer))

	// 1) Dial gRPC client
	conn, err := grpc.NewClient(s.grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.T(), err)
	defer func() { _ = conn.Close() }()

	client := customerv1.NewCustomerAuthServiceClient(conn)

	// 2) Build request
	req := &customerv1.LoginRequest{Phone: customer.Phone, Password: "P@ssw0rd!"}

	// 3) Call Login
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()
	res, err := client.Login(ctx, req)

	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.GetChallengeId())

	ok := s.waitForIMAPSent("Your OTP Code")
	require.True(s.T(), ok)
}

func (s *LoginE2ESuite) waitForIMAPSent(subject string) bool {
	require.NotEmpty(s.T(), s.tCfg.ImapHost)
	require.NotEmpty(s.T(), s.tCfg.ImapPort)
	require.NotEmpty(s.T(), s.tCfg.ImapUsername)
	require.NotEmpty(s.T(), s.tCfg.ImapPassword)

	deadline := time.Now().Add(10 * time.Second)
	addr := net.JoinHostPort(s.tCfg.ImapHost, s.tCfg.ImapPort)

	c, err := imapclient.DialTLS(addr, &tls.Config{ServerName: s.tCfg.ImapHost})
	require.NoError(s.T(), err)

	err = c.Login(s.tCfg.ImapUsername, s.tCfg.ImapPassword)
	require.NoError(s.T(), err)
	defer func() { _ = c.Logout() }()

	_, err = c.Select("INBOX", false)
	require.NoError(s.T(), err)

	for time.Now().Before(deadline) {
		crit := imap.NewSearchCriteria()
		crit.Since = time.Now().Add(-1 * time.Minute)
		crit.Header.Add("Subject", subject)
		crit.Header.Add("To", s.tCfg.ImapUsername)

		ids, err := c.Search(crit)
		s.T().Log("ids", ids)
		if err == nil && len(ids) > 0 {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}

	return false
}

func Test_LoginE2E(t *testing.T) {
	suite.Run(t, new(LoginE2ESuite))
}
