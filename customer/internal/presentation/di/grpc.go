package di

import (
	"context"
	"customer/internal/infrastructure/logger"
	customerv1 "customer/internal/presentation/grpc"
	"customer/internal/presentation/grpc/handler"
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
		customerv1.NewConfig,

		// Handlers
		fx.Annotate(
			handler.NewCustomerAuthServiceHandler,
			fx.As(new(customerv1.CustomerAuthServiceServer)),
		),

		// GRPC server
		newGRPCServer,
	),

	// Lifecycle
	fx.Invoke(setupGRPCLifecycle),
)

func newGRPCServer(customerAuthHandler customerv1.CustomerAuthServiceServer) *grpc.Server {
	server := grpc.NewServer()
	customerv1.RegisterCustomerAuthServiceServer(server, customerAuthHandler)
	reflection.Register(server)
	return server
}

func setupGRPCLifecycle(lc fx.Lifecycle, cfg *customerv1.Config, server *grpc.Server, logger logger.Logger) {
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
