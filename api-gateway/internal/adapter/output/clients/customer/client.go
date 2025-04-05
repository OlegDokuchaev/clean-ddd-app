package customer

import (
	customerGRPC "api-gateway/gen/customer/v1"
	"api-gateway/internal/adapter/output/clients/response"
	customerDto "api-gateway/internal/domain/dtos/customer"
	customerClient "api-gateway/internal/port/output/clients/customer"
	"context"

	"github.com/google/uuid"
)

type ClientImpl struct {
	client customerGRPC.CustomerAuthServiceClient
}

func NewClient(client customerGRPC.CustomerAuthServiceClient) customerClient.Client {
	return &ClientImpl{
		client: client,
	}
}

func (c *ClientImpl) Register(ctx context.Context, data customerDto.RegisterDto) (uuid.UUID, error) {
	request := toRegisterRequest(data)

	resp, err := c.client.Register(ctx, request)
	if err != nil {
		return uuid.Nil, response.ParseGRPCError(err)
	}

	customerID, err := response.ToUUID(resp.CustomerId)
	if err != nil {
		return uuid.Nil, err
	}

	return customerID, nil
}

func (c *ClientImpl) Login(ctx context.Context, data customerDto.LoginDto) (string, error) {
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

	customerID, err := response.ToUUID(resp.CustomerId)
	if err != nil {
		return uuid.Nil, err
	}

	return customerID, nil
}

var _ customerClient.Client = (*ClientImpl)(nil)
