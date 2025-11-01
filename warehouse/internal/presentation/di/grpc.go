package di

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"

	"warehouse/internal/infrastructure/logger"
	warehousev1 "warehouse/internal/presentation/grpc"
	"warehouse/internal/presentation/grpc/handlers"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCModule aggregates all components for the gRPC server.
var GRPCModule = fx.Options(
	fx.Provide(
		// Config
		warehousev1.NewConfig,

		// RPC method handlers
		fx.Annotate(
			handlers.NewItemServiceHandler,
			fx.As(new(warehousev1.ItemServiceServer)),
		),
		fx.Annotate(
			handlers.NewProductServiceHandler,
			fx.As(new(warehousev1.ProductServiceServer)),
		),
		fx.Annotate(
			handlers.NewProductImageServiceHandler,
			fx.As(new(warehousev1.ProductImageServiceServer)),
		),

		// gRPC server
		newGRPCServer,
	),

	// Lifecycle
	fx.Invoke(setupGRPCLifecycle),
)

func newGRPCServer(
	itemHandler warehousev1.ItemServiceServer,
	productHandler warehousev1.ProductServiceServer,
	productImageHandler warehousev1.ProductImageServiceServer,
	logger logger.Logger,
) *grpc.Server {
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents),
		)),

		grpc.UnaryInterceptor(warehousev1.LoggingInterceptor(logger)),
		grpc.StreamInterceptor(warehousev1.StreamLoggingInterceptor(logger)),
	)

	warehousev1.RegisterItemServiceServer(server, itemHandler)
	warehousev1.RegisterProductServiceServer(server, productHandler)
	warehousev1.RegisterProductImageServiceServer(server, productImageHandler)
	reflection.Register(server)

	return server
}

func setupGRPCLifecycle(lc fx.Lifecycle, cfg *warehousev1.Config, server *grpc.Server, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%s", cfg.Port)
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return fmt.Errorf("failed to listen on %s: %w", address, err)
			}

			go func() {
				if err := server.Serve(listener); err != nil {
					logger.Printf("error starting gRPC server: %v", err)
				}
			}()

			logger.Printf("gRPC server started on port: %s", cfg.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Println("gRPC server stopping...")
			server.GracefulStop()
			logger.Println("gRPC server stopped")
			return nil
		},
	})
}
