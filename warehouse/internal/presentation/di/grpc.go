package di

import (
	"context"
	"fmt"
	"log"
	"net"

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

		// gRPC server
		newGRPCServer,
	),

	// Lifecycle
	fx.Invoke(setupGRPCLifecycle),
)

func newGRPCServer(
	itemHandler warehousev1.ItemServiceServer,
	productHandler warehousev1.ProductServiceServer,
) *grpc.Server {
	server := grpc.NewServer()
	warehousev1.RegisterItemServiceServer(server, itemHandler)
	warehousev1.RegisterProductServiceServer(server, productHandler)
	reflection.Register(server)
	return server
}

func setupGRPCLifecycle(lc fx.Lifecycle, cfg *warehousev1.Config, server *grpc.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%s", cfg.Port)
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return fmt.Errorf("failed to listen on %s: %w", address, err)
			}

			go func() {
				if err := server.Serve(listener); err != nil {
					log.Printf("error starting gRPC server: %v", err)
				}
			}()

			log.Printf("gRPC server started on port: %s", cfg.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("gRPC server stopping...")
			server.GracefulStop()
			log.Println("gRPC server stopped")
			return nil
		},
	})
}
