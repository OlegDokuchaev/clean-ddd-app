package di

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"warehouse/internal/infrastructure/messaging"

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

		// Message writers
		fx.Annotate(
			messaging.NewWarehouseCmdResWriter,
			fx.ResultTags(`name:"warehouseCmdResWriter"`),
		),
	),

	// Kafka resources lifecycle management
	fx.Invoke(setupMessagingLifecycle),
)

func setupMessagingLifecycle(in struct {
	fx.In

	Lifecycle fx.Lifecycle

	// Readers
	WarehouseCmdReader *kafka.Reader `name:"warehouseCmdReader"`

	// Writers
	WarehouseCmdResWriter *kafka.Writer `name:"warehouseCmdResWriter"`
}) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Kafka resources ready for use")
			return nil
		},

		OnStop: func(ctx context.Context) error {
			log.Println("Closing Kafka resources...")
			var hasErrors bool

			// Close readers
			if err := closeReader("warehouse command reader", in.WarehouseCmdReader); err != nil {
				hasErrors = true
			}

			// Close writers
			if err := closeWriter("warehouse command result writer", in.WarehouseCmdResWriter); err != nil {
				hasErrors = true
			}

			if hasErrors {
				return fmt.Errorf("errors occurred while closing Kafka resources")
			}
			log.Println("All Kafka resources successfully closed")
			return nil
		},
	})
}

func closeReader(name string, reader *kafka.Reader) error {
	if reader == nil {
		return nil
	}

	if err := reader.Close(); err != nil {
		log.Printf("Error closing %s: %v", name, err)
		return err
	}

	return nil
}

func closeWriter(name string, writer *kafka.Writer) error {
	if writer == nil {
		return nil
	}

	if err := writer.Close(); err != nil {
		log.Printf("Error closing %s: %v", name, err)
		return err
	}

	return nil
}
