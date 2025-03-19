package di

import (
	"context"
	"fmt"
	"net"
	orderv1 "order/internal/presentation/grpc"
	"order/internal/presentation/grpc/handler"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var GRPCModule = fx.Provide(
	orderv1.NewConfig,
	fx.Annotate(
		handler.NewOrderServiceHandler,
		fx.As(new(orderv1.OrderServiceServer)),
	),
	newGRPCServer,
)

func newGRPCServer(lc fx.Lifecycle, cfg *orderv1.Config, orderHandler orderv1.OrderServiceServer) *grpc.Server {
	server := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(server, orderHandler)
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
