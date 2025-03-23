package di

import (
	"context"
	courierv1 "courier/internal/presentation/grpc"
	"courier/internal/presentation/grpc/handler"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var GRPCModule = fx.Provide(
	// Configuration
	courierv1.NewConfig,

	// Handlers
	fx.Annotate(
		handler.NewCourierAuthServiceHandler,
		fx.As(new(courierv1.CourierAuthServiceServer)),
	),

	// GRPC server
	newGRPCServer,
)

func newGRPCServer(
	lc fx.Lifecycle,
	cfg *courierv1.Config,
	courierAuthHandler courierv1.CourierAuthServiceServer,
) *grpc.Server {
	server := grpc.NewServer()
	courierv1.RegisterCourierAuthServiceServer(server, courierAuthHandler)
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
