package main

import (
	"context"
	appDI "customer/internal/application/di"
	infraDI "customer/internal/infrastructure/di"
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
		infraDI.DatabaseModule,
		infraDI.RepositoryModule,
		infraDI.TokenManagerModule,

		// Application modules
		appDI.UseCaseModule,

		// Presentation modules
		presentationDI.GRPCModule,

		// Add logging for application startup and shutdown
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					log.Println("Starting Customer service...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Println("Shutting down Customer service...")
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
