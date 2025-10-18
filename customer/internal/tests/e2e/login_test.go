//go:build e2e

package e2e

import (
	"context"
	"crypto/tls"
	"net"
	"testing"
	"time"

	appDI "customer/internal/application/di"
	"customer/internal/infrastructure/auth"
	"customer/internal/infrastructure/db/migrations"
	infraDI "customer/internal/infrastructure/di"
	"customer/internal/infrastructure/logger"
	customerRepository "customer/internal/infrastructure/repository/customer"
	presentationDI "customer/internal/presentation/di"
	customerv1 "customer/internal/presentation/grpc"
	"customer/internal/tests/testutils"
	"customer/internal/tests/testutils/mothers"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ---- Global test fixtures (equivalent to testify's suite fields) ----
var (
	testCtx context.Context

	app     *fx.App
	tCfg    *testutils.Config
	grpcURL string

	db    *testutils.TestDB
	redis *testutils.TestRedis
)

// ---- Helpers (equivalent to clear/waitForIMAPSent) ----

// clearAll wipes test DB and Redis between specs and during suite teardown.
// It mirrors the behavior of your previous clear() helper.
func clearAll(ctx context.Context) {
	if db != nil {
		Expect(db.Clear(ctx)).To(Succeed())
	}
	if redis != nil {
		Expect(redis.Clear(ctx)).To(Succeed())
	}
}

// waitForIMAPSent polls the IMAP inbox for a message with the given subject
// addressed to tCfg.ImapUsername. Returns true if found before the deadline.
// This function keeps the original behavior and timing from your suite.
func waitForIMAPSent(subject string) bool {
	Expect(tCfg.ImapHost).NotTo(BeEmpty())
	Expect(tCfg.ImapPort).NotTo(BeEmpty())
	Expect(tCfg.ImapUsername).NotTo(BeEmpty())
	Expect(tCfg.ImapPassword).NotTo(BeEmpty())

	addr := net.JoinHostPort(tCfg.ImapHost, tCfg.ImapPort)

	// Use TLS with SNI set to the IMAP host.
	c, err := imapclient.DialTLS(addr, &tls.Config{ServerName: tCfg.ImapHost})
	if err != nil {
		return false
	}
	defer func() { _ = c.Logout() }()

	if err := c.Login(tCfg.ImapUsername, tCfg.ImapPassword); err != nil {
		return false
	}
	if _, err := c.Select("INBOX", false); err != nil {
		return false
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		crit := imap.NewSearchCriteria()
		crit.Since = time.Now().Add(-1 * time.Minute)
		crit.Header.Add("Subject", subject)
		crit.Header.Add("To", tCfg.ImapUsername)

		ids, err := c.Search(crit)
		if err == nil && len(ids) > 0 {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

// ---- Suite lifecycle (Ginkgo equivalents of SetupSuite/TearDownSuite/AfterTest) ----

var _ = BeforeSuite(func(ctx SpecContext) {
	testCtx = context.Background()

	// 1) Load configs
	var err error
	tCfg, err = testutils.NewConfig()
	Expect(err).NotTo(HaveOccurred())

	mCfg, err := migrations.NewConfig()
	Expect(err).NotTo(HaveOccurred())

	// 2) Test db
	db, err = testutils.NewTestDB(testCtx, tCfg, mCfg)
	Expect(err).NotTo(HaveOccurred())

	// 3) Test redis
	redis, err = testutils.NewTestRedis(testCtx, tCfg)
	Expect(err).NotTo(HaveOccurred())

	clearAll(testCtx)

	// 4) GRPC
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).NotTo(HaveOccurred())
	_, grpcPortStr, err := net.SplitHostPort(ln.Addr().String())
	Expect(err).NotTo(HaveOccurred())
	_ = ln.Close()
	grpcURL = net.JoinHostPort("127.0.0.1", grpcPortStr)

	grpcCfg := &customerv1.Config{Port: grpcPortStr}

	// 5) Logrus
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	// 5) Auth config
	aCfg := &auth.Config{
		SigningKey:       "test-signing-key",
		AccessTTL:        15 * time.Minute,
		ResetTTL:         1 * time.Hour,
		LockoutMaxFailed: 5,
		LockoutLockFor:   30 * time.Minute,
		OtpTTL:           5 * time.Minute,
		OtpMaxAttempts:   3,
	}

	// 7) FX app
	app = fx.New(
		infraDI.LoggerModule,
		infraDI.DatabaseModule,
		infraDI.RepositoryModule,
		infraDI.TokenManagerModule,
		infraDI.MailSenderModule,
		infraDI.OtpStoreModule,
		infraDI.AuthPoliciesModule,
		appDI.UseCaseModule,
		presentationDI.GRPCModule,

		// Replace infra with test instances/configs.
		fx.Replace(logrus.New()),
		fx.Replace(db.Cfg),
		fx.Replace(db.DB),
		fx.Replace(redis.Cfg),
		fx.Replace(redis.Client),
		fx.Replace(grpcCfg),
		fx.Replace(aCfg),

		// No-op lifecycle hooks to maintain parity with the original code.
		fx.Invoke(func(lc fx.Lifecycle, l logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error { return nil },
				OnStop:  func(context.Context) error { return nil },
			})
		}),
	)

	startCtx, cancel := context.WithTimeout(testCtx, time.Minute)
	defer cancel()
	Expect(app.Start(startCtx)).To(Succeed())

	// 8) Wait until gRPC is ready
	Eventually(func() error {
		c, err := net.DialTimeout("tcp", grpcURL, 2*time.Second)
		if err != nil {
			return err
		}
		_ = c.Close()
		return nil
	}).WithTimeout(10 * time.Second).WithPolling(200 * time.Millisecond).Should(Succeed())
})

var _ = AfterEach(func(ctx SpecContext) {
	clearAll(testCtx)
})

var _ = AfterSuite(func() {
	if app != nil {
		ctx, cancel := context.WithTimeout(testCtx, 20*time.Second)
		_ = app.Stop(ctx)
		cancel()
	}
	clearAll(testCtx)
})

var _ = Describe("Login E2E", func() {
	It("sends a real email and returns a challenge", func(ctx SpecContext) {
		customer := mothers.CustomerWithPassword(tCfg.ImapPassword)
		customer.Email = tCfg.ImapUsername

		r := customerRepository.New(db.DB)
		Expect(r.Create(testCtx, customer)).To(Succeed())

		conn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).NotTo(HaveOccurred())
		defer func() { _ = conn.Close() }()

		client := customerv1.NewCustomerAuthServiceClient(conn)

		req := &customerv1.LoginRequest{
			Phone:    customer.Phone,
			Password: tCfg.ImapPassword,
		}

		callCtx, cancel := context.WithTimeout(testCtx, 30*time.Second)
		defer cancel()

		res, err := client.Login(callCtx, req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.GetChallengeId()).NotTo(BeEmpty())

		Expect(waitForIMAPSent("Your OTP Code")).To(BeTrue())
	}, NodeTimeout(60*time.Second))
})

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Customer Auth E2E Suite")
}
