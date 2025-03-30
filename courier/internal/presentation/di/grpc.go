package di

import (
	"context"
	"courier/internal/infrastructure/logger"
	courierv1 "courier/internal/presentation/grpc"
	"courier/internal/presentation/grpc/handler"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCModule aggregates all components for the gRPC server.
var GRPCModule = fx.Options(
	fx.Provide(
		// Configuration
		courierv1.NewConfig,

		// Handlers
		fx.Annotate(
			handler.NewCourierAuthServiceHandler,
			fx.As(new(courierv1.CourierAuthServiceServer)),
		),

		// GRPC server
		newGRPCServer,
	),

	// Lifecycle
	fx.Invoke(setupGRPCLifecycle),
)

func newGRPCServer(courierAuthHandler courierv1.CourierAuthServiceServer) *grpc.Server {
	server := grpc.NewServer()
	courierv1.RegisterCourierAuthServiceServer(server, courierAuthHandler)
	reflection.Register(server)
	return server
}

func setupGRPCLifecycle(lc fx.Lifecycle, cfg *courierv1.Config, server *grpc.Server, logger logger.Logger) {
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
