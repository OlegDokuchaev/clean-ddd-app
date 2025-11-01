package di

import (
	"context"
	"fmt"
	"net"
	"order/internal/infrastructure/logger"
	orderv1 "order/internal/presentation/grpc"
	"order/internal/presentation/grpc/handler"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCModule aggregates all components for the gRPC server.
var GRPCModule = fx.Options(
	fx.Provide(
		// Configuration
		orderv1.NewConfig,

		// Handlers
		fx.Annotate(
			handler.NewOrderServiceHandler,
			fx.As(new(orderv1.OrderServiceServer)),
		),

		// GRPC server
		newGRPCServer,
	),

	// Lifecycle
	fx.Invoke(setupGRPCLifecycle),
)

func newGRPCServer(orderHandler orderv1.OrderServiceServer, logger logger.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents),
		)),

		grpc.UnaryInterceptor(orderv1.LoggingInterceptor(logger)),
		grpc.StreamInterceptor(orderv1.StreamLoggingInterceptor(logger)),
	)

	orderv1.RegisterOrderServiceServer(server, orderHandler)
	reflection.Register(server)
	return server
}

func setupGRPCLifecycle(lc fx.Lifecycle, cfg *orderv1.Config, server *grpc.Server, logger logger.Logger) {
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
