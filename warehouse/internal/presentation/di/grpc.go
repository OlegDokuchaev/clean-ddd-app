package di

import (
	"context"
	"fmt"
	"net"
	warehousev1 "warehouse/internal/presentation/grpc"
	"warehouse/internal/presentation/grpc/handlers"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var GRPCModule = fx.Provide(
	// Configuration
	warehousev1.NewConfig,

	// Handlers
	fx.Annotate(
		handlers.NewItemServiceHandler,
		fx.As(new(warehousev1.ItemServiceServer)),
	),
	fx.Annotate(
		handlers.NewProductServiceHandler,
		fx.As(new(warehousev1.ProductServiceServer)),
	),

	// GRPC server
	newGRPCServer,
)

func newGRPCServer(
	lc fx.Lifecycle,
	cfg *warehousev1.Config,
	itemHandler warehousev1.ItemServiceServer,
	productHandler warehousev1.ProductServiceServer,
) *grpc.Server {
	server := grpc.NewServer()
	warehousev1.RegisterItemServiceServer(server, itemHandler)
	warehousev1.RegisterProductServiceServer(server, productHandler)
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
