package courier

import (
	courierGRPC "api-gateway/gen/courier/v1"
	"api-gateway/internal/adapter/output/clients/response"
	courierDto "api-gateway/internal/domain/dtos/courier"
	courierClient "api-gateway/internal/port/output/clients/courier"
	"context"

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
	request := toLoginRequest(data)

	resp, err := c.client.Login(ctx, request)
	if err != nil {
		return "", response.ParseGRPCError(err)
	}

	return resp.Token, nil
}

func (c *ClientImpl) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
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
