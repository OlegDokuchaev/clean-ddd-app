package di

import (
	"context"
	"fmt"
	"log"
	"net"
	orderv1 "order/internal/presentation/grpc"
	"order/internal/presentation/grpc/handler"

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

func newGRPCServer(orderHandler orderv1.OrderServiceServer) *grpc.Server {
	server := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(server, orderHandler)
	reflection.Register(server)
	return server
}

func setupGRPCLifecycle(lc fx.Lifecycle, cfg *orderv1.Config, server *grpc.Server) {
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
