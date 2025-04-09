package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type GRPCClients struct {
	Product warehouseGRPC.ProductServiceClient
	Item    warehouseGRPC.ItemServiceClient
}

func newConnection(config *Config) (*grpc.ClientConn, error) {
	timeout := time.Duration(config.TimeoutSeconds) * time.Second
	target := "passthrough:///" + config.Address

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: timeout}),
	)
	if err != nil {
		return nil, err
	}

	conn.Connect()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	state := conn.GetState()
	for state != connectivity.Ready {
		if !conn.WaitForStateChange(ctx, state) {
			return nil, ctx.Err()
		}
		state = conn.GetState()
	}

	return conn, nil
}

func NewGRPCClient(config *Config) (*GRPCClients, error) {
	conn, err := newConnection(config)
	if err != nil {
		return nil, err
	}
	return &GRPCClients{
		Product: warehouseGRPC.NewProductServiceClient(conn),
		Item:    warehouseGRPC.NewItemServiceClient(conn),
	}, nil
}
