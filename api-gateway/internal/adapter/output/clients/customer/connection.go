package customer

import (
	customerGRPC "api-gateway/gen/customer/v1"
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func newConnection(config *Config) (*grpc.ClientConn, error) {
	timeout := time.Duration(config.TimeoutSeconds) * time.Second
	target := "passthrough:///" + config.Address

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: timeout}),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
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

func NewGRPCClient(config *Config) (customerGRPC.CustomerAuthServiceClient, error) {
	conn, err := newConnection(config)
	if err != nil {
		return nil, err
	}
	return customerGRPC.NewCustomerAuthServiceClient(conn), nil
}
