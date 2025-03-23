package di

import (
	"context"
	customerv1 "customer/internal/presentation/grpc"
	"customer/internal/presentation/grpc/handler"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var GRPCModule = fx.Provide(
	// Configuration
	customerv1.NewConfig,

	// Handlers
	fx.Annotate(
		handler.NewCustomerAuthServiceHandler,
		fx.As(new(customerv1.CustomerAuthServiceServer)),
	),

	// GRPC server
	newGRPCServer,
)

func newGRPCServer(
	lc fx.Lifecycle,
	cfg *customerv1.Config,
	customerAuthHandler customerv1.CustomerAuthServiceServer,
) *grpc.Server {
	server := grpc.NewServer()
	customerv1.RegisterCustomerAuthServiceServer(server, customerAuthHandler)
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
			if err != nil {
				return fmt.Errorf("failed to start gRPC server: %w", err)
			}

			go func() {
				if err := server.Serve(listener); err != nil {
					fmt.Printf("error starting gRPC server: %v\n", err)
				}
			}()

			fmt.Printf("gRPC server started on port: %s\n", cfg.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			fmt.Println("gRPC server stopped")
			return nil
		},
	})

	return server
}
