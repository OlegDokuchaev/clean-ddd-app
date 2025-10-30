package main

import (
	"context"
	appDI "customer/internal/application/di"
	infraDI "customer/internal/infrastructure/di"
	"customer/internal/infrastructure/logger"
	presentationDI "customer/internal/presentation/di"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Infrastructure modules
		infraDI.LoggerModule,
		infraDI.DatabaseModule,
		infraDI.RepositoryModule,
		infraDI.TokenManagerModule,
		infraDI.MailSenderModule,
		infraDI.OtpStoreModule,
		infraDI.AuthPoliciesModule,
		infraDI.TelemetryModule,

		// Application modules
		appDI.UseCaseModule,

		// Presentation modules
		presentationDI.GRPCModule,
		presentationDI.TelemetryModule,

		// Add logging for application startup and shutdown
		fx.Invoke(func(lc fx.Lifecycle, logger logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Println("Starting Customer service...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Println("Shutting down Customer service...")
					return nil
				},
			})
		}),
	)

	// Setting up proper application termination on signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capturing termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Starting the application
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Waiting for termination signal
	sig := <-sigCh
	log.Printf("Received signal: %v", sig)

	// Graceful application shutdown
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatalf("Failed to stop application: %v", err)
	}
}
