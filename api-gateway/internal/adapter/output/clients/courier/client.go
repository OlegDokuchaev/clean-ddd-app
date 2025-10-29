package courier

import (
	courierGRPC "api-gateway/gen/courier/v1"
	"api-gateway/internal/adapter/output/clients/response"
	courierDto "api-gateway/internal/domain/dtos/courier"
	courierClient "api-gateway/internal/port/output/clients/courier"
	"context"
	"go.opentelemetry.io/otel"

	"github.com/google/uuid"
)

type ClientImpl struct {
	client courierGRPC.CourierAuthServiceClient
}

func NewClient(client courierGRPC.CourierAuthServiceClient) courierClient.Client {
	return &ClientImpl{
		client: client,
	}
}

func (c *ClientImpl) Register(ctx context.Context, data courierDto.RegisterDto) (uuid.UUID, error) {
	ctx, span := otel.Tracer("api-gateway.courier").Start(ctx, "Register")
	defer span.End()

	request := toRegisterRequest(data)

	resp, err := c.client.Register(ctx, request)
	if err != nil {
		return uuid.Nil, response.ParseGRPCError(err)
	}

	courierID, err := response.ToUUID(resp.CourierId)
	if err != nil {
		return uuid.Nil, err
	}

	return courierID, nil
}

func (c *ClientImpl) Login(ctx context.Context, data courierDto.LoginDto) (string, error) {
	ctx, span := otel.Tracer("api-gateway.courier").Start(ctx, "Login")
	defer span.End()

	request := toLoginRequest(data)

	resp, err := c.client.Login(ctx, request)
	if err != nil {
		return "", response.ParseGRPCError(err)
	}

	return resp.Token, nil
}

func (c *ClientImpl) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	ctx, span := otel.Tracer("api-gateway.courier").Start(ctx, "Authenticate")
	defer span.End()

	request := toAuthenticateRequest(token)

	resp, err := c.client.Authenticate(ctx, request)
	if err != nil {
		return uuid.Nil, response.ParseGRPCError(err)
	}

	courierID, err := response.ToUUID(resp.CourierId)
	if err != nil {
		return uuid.Nil, err
	}

	return courierID, nil
}

var _ courierClient.Client = (*ClientImpl)(nil)
