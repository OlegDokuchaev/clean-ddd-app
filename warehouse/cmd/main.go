package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	appDI "warehouse/internal/application/di"
	infraDI "warehouse/internal/infrastructure/di"
	"warehouse/internal/infrastructure/logger"
	presentationDI "warehouse/internal/presentation/di"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Infrastructure modules
		infraDI.LoggerModule,
		infraDI.MessagingModule,
		infraDI.DatabaseModule,
		infraDI.RepositoryModule,
		infraDI.PublisherModule,
		infraDI.OutboxProcessorModule,
		infraDI.UowModule,

		// Application modules
		appDI.UseCaseModule,

		// Presentation modules
		presentationDI.GRPCModule,
		presentationDI.CommandConsumerModule,
		presentationDI.EventsModule,

		// Add logging for application startup and shutdown
		fx.Invoke(func(lc fx.Lifecycle, logger logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Println("Starting Warehouse service...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Println("Shutting down Warehouse service...")
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
