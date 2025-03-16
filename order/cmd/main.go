package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"

	appDI "order/internal/application/di"
	infraDI "order/internal/infrastructure/di"
	presentationDI "order/internal/presentation/di"
)

func main() {
	app := fx.New(
		infraDI.RepositoryModule,
		infraDI.PublisherModule,
		appDI.UseCaseModule,
		appDI.SagaModule,
		presentationDI.GRPCModule,
		presentationDI.CommandsModule,
		presentationDI.SagaModule,

		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					log.Println("Starting Order service...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Println("Shutting down Order service...")
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
