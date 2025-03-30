package di

import (
	"context"
	"fmt"
	"warehouse/internal/infrastructure/logger"
	"warehouse/internal/infrastructure/messaging"

	"github.com/segmentio/kafka-go"

	"go.uber.org/fx"
)

var MessagingModule = fx.Options(
	fx.Provide(
		// Configuration
		messaging.NewConfig,

		// Message readers
		fx.Annotate(
			messaging.NewWarehouseCmdReader,
			fx.ResultTags(`name:"warehouseCmdReader"`),
		),

		fx.Annotate(
			messaging.NewProductEventReader,
			fx.ResultTags(`name:"productEventReader"`),
		),

		// Message writers
		fx.Annotate(
			messaging.NewWarehouseCmdResWriter,
			fx.ResultTags(`name:"warehouseCmdResWriter"`),
		),

		fx.Annotate(
			messaging.NewProductEventWriter,
			fx.ResultTags(`name:"productEventWriter"`),
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
	WarehouseCmdReader *kafka.Reader `name:"warehouseCmdReader"`
	ProductEventReader *kafka.Reader `name:"productEventReader"`

	// Writers
	WarehouseCmdResWriter *kafka.Writer `name:"warehouseCmdResWriter"`
	ProductEventWriter    *kafka.Writer `name:"productEventWriter"`
}) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			in.Logger.Println("Kafka resources ready for use")
			return nil
		},

		OnStop: func(ctx context.Context) error {
			in.Logger.Println("Closing Kafka resources...")
			var hasErrors bool

			// Close readers
			if err := closeReader("warehouse command reader", in.WarehouseCmdReader, in.Logger); err != nil {
				hasErrors = true
			}

			if err := closeReader("product event reader", in.ProductEventReader, in.Logger); err != nil {
				hasErrors = true
			}

			// Close writers
			if err := closeWriter("warehouse command result writer", in.WarehouseCmdResWriter, in.Logger); err != nil {
				hasErrors = true
			}

			if err := closeWriter("product event writer", in.ProductEventWriter, in.Logger); err != nil {
				hasErrors = true
			}

			if hasErrors {
				return fmt.Errorf("errors occurred while closing Kafka resources")
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
