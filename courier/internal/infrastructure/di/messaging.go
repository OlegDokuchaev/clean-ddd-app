package di

import (
	"context"
	"courier/internal/infrastructure/logger"
	"courier/internal/infrastructure/messaging"
	"errors"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

var MessagingModule = fx.Options(
	fx.Provide(
		// Configuration
		messaging.NewConfig,

		// Readers
		fx.Annotate(
			messaging.NewCourierCmdReader,
			fx.ResultTags(`name:"courierCmdReader"`),
		),

		// Writers
		fx.Annotate(
			messaging.NewCourierCmdResWriter,
			fx.ResultTags(`name:"courierCmdResWriter"`),
		),
	),

	// Kafka resources lifecycle management
	fx.Invoke(setupMessagingLifecycle),
)

func setupMessagingLifecycle(in struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    logger.Logger

	// Readers
	CourierCmdReader *kafka.Reader `name:"courierCmdReader"`

	// Writers
	CourierCmdResWriter *kafka.Writer `name:"courierCmdResWriter"`
}) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			in.Logger.Println("Kafka resources ready for use")
			return nil
		},

		OnStop: func(ctx context.Context) error {
			in.Logger.Println("Closing Kafka resources...")
			var errs []error

			// Close readers
			if err := closeReader("courier command reader", in.CourierCmdReader, in.Logger); err != nil {
				errs = append(errs, err)
			}

			// Close writers
			if err := closeWriter("courier command result writer", in.CourierCmdResWriter, in.Logger); err != nil {
				errs = append(errs, err)
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			in.Logger.Println("All Kafka resources successfully closed")
			return nil
		},
	})
}

func closeReader(name string, reader *kafka.Reader, logger logger.Logger) error {
	if reader == nil {
		return nil
	}

	if err := reader.Close(); err != nil {
		logger.Printf("Error closing %s: %v", name, err)
		return err
	}

	return nil
}

func closeWriter(name string, writer *kafka.Writer, logger logger.Logger) error {
	if writer == nil {
		return nil
	}

	if err := writer.Close(); err != nil {
		logger.Printf("Error closing %s: %v", name, err)
		return err
	}

	return nil
}
