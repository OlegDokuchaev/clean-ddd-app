package di

import (
	"context"
	"fmt"
	"order/internal/infrastructure/logger"
	"order/internal/infrastructure/messaging"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

var MessagingModule = fx.Options(
	fx.Provide(
		// General Kafka configuration
		messaging.NewConfig,

		// Message readers
		fx.Annotate(
			messaging.NewOrderCommandReader,
			fx.ResultTags(`name:"orderCommandReader"`),
		),
		fx.Annotate(
			messaging.NewWarehouseCommandResultReader,
			fx.ResultTags(`name:"warehouseCommandResultReader"`),
		),
		fx.Annotate(
			messaging.NewCourierCommandResultReader,
			fx.ResultTags(`name:"courierCommandResultReader"`),
		),

		// Message writers
		fx.Annotate(
			messaging.NewOrderCommandWriter,
			fx.ResultTags(`name:"orderCommandWriter"`),
		),
		fx.Annotate(
			messaging.NewWarehouseCommandWriter,
			fx.ResultTags(`name:"warehouseCommandWriter"`),
		),
		fx.Annotate(
			messaging.NewCourierCommandWriter,
			fx.ResultTags(`name:"courierCommandWriter"`),
		),
		fx.Annotate(
			messaging.NewOrderCommandResWriter,
			fx.ResultTags(`name:"orderCommandResWriter"`),
		),
	),

	// Kafka resources lifecycle management
	fx.Invoke(setupMessagingLifecycle),
)

// setupMessagingLifecycle configures centralized closing of Kafka resources
func setupMessagingLifecycle(in struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    logger.Logger

	// Readers
	OrderCommandReader        *kafka.Reader `name:"orderCommandReader"`
	WarehouseCommandResReader *kafka.Reader `name:"warehouseCommandResultReader"`
	CourierCommandResReader   *kafka.Reader `name:"courierCommandResultReader"`

	// Writers
	OrderCommandWriter     *kafka.Writer `name:"orderCommandWriter"`
	WarehouseCommandWriter *kafka.Writer `name:"warehouseCommandWriter"`
	CourierCommandWriter   *kafka.Writer `name:"courierCommandWriter"`
	OrderCommandResWriter  *kafka.Writer `name:"orderCommandResWriter"`
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
			if err := closeReader("order command reader", in.OrderCommandReader, in.Logger); err != nil {
				hasErrors = true
			}
			if err := closeReader("warehouse command result reader", in.WarehouseCommandResReader, in.Logger); err != nil {
				hasErrors = true
			}
			if err := closeReader("courier command result reader", in.CourierCommandResReader, in.Logger); err != nil {
				hasErrors = true
			}

			// Close writers
			if err := closeWriter("order command writer", in.OrderCommandWriter, in.Logger); err != nil {
				hasErrors = true
			}
			if err := closeWriter("warehouse command writer", in.WarehouseCommandWriter, in.Logger); err != nil {
				hasErrors = true
			}
			if err := closeWriter("courier command writer", in.CourierCommandWriter, in.Logger); err != nil {
				hasErrors = true
			}
			if err := closeWriter("order command response writer", in.OrderCommandResWriter, in.Logger); err != nil {
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
