package customer

import (
	customerGRPC "api-gateway/gen/customer/v1"
	"api-gateway/internal/adapter/output/clients/response"
	customerDto "api-gateway/internal/domain/dtos/customer"
	customerClient "api-gateway/internal/port/output/clients/customer"
	"context"
	"go.opentelemetry.io/otel"

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
	ctx, span := otel.Tracer("api-gateway.customer").Start(ctx, "Register")
	defer span.End()

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
	ctx, span := otel.Tracer("api-gateway.customer").Start(ctx, "Login")
	defer span.End()

	request := toLoginRequest(data)

	resp, err := c.client.Login(ctx, request)
	if err != nil {
		return "", response.ParseGRPCError(err)
	}

	return resp.ChallengeId, nil
}

func (c *ClientImpl) VerifyOtp(ctx context.Context, data customerDto.VerifyOtpDto) (string, error) {
	ctx, span := otel.Tracer("api-gateway.customer").Start(ctx, "VerifyOtp")
	defer span.End()

	request := toVerifyOtpRequest(data)

	resp, err := c.client.VerifyOtp(ctx, request)
	if err != nil {
		return "", response.ParseGRPCError(err)
	}

	return resp.Token, nil
}

func (c *ClientImpl) RequestPasswordReset(ctx context.Context, email string) error {
	ctx, span := otel.Tracer("api-gateway.customer").Start(ctx, "RequestPasswordReset")
	defer span.End()

	request := toRequestPasswordResetRequest(email)

	_, err := c.client.RequestPasswordReset(ctx, request)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) CompletePasswordReset(ctx context.Context, token string, newPassword string) error {
	ctx, span := otel.Tracer("api-gateway.customer").Start(ctx, "CompletePasswordReset")
	defer span.End()

	request := toCompletePasswordResetRequest(token, newPassword)

	_, err := c.client.CompletePasswordReset(ctx, request)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	ctx, span := otel.Tracer("api-gateway.customer").Start(ctx, "Authenticate")
	defer span.End()

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
