package handler

import (
	"context"
	courierAuthApplication "courier/internal/application/courier/auth"
	courierv1 "courier/internal/presentation/grpc"
	"courier/internal/presentation/grpc/request"
	"courier/internal/presentation/grpc/response"
)

type CourierAuthServiceHandler struct {
	courierv1.UnimplementedCourierAuthServiceServer

	usecase courierAuthApplication.UseCase
}

func NewCourierAuthServiceHandler(usecase courierAuthApplication.UseCase) *CourierAuthServiceHandler {
	return &CourierAuthServiceHandler{
		usecase: usecase,
	}
}

func (h *CourierAuthServiceHandler) Register(
	ctx context.Context,
	req *courierv1.RegisterRequest,
) (*courierv1.RegisterResponse, error) {
	data, err := request.ToRegisterDto(req)
	if err != nil {
		return nil, err
	}

	courierID, err := h.usecase.Register(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToRegisterResponse(courierID), nil
}

func (h *CourierAuthServiceHandler) Login(
	ctx context.Context,
	req *courierv1.LoginRequest,
) (*courierv1.LoginResponse, error) {
	data, err := request.ToLoginDto(req)
	if err != nil {
		return nil, err
	}

	token, err := h.usecase.Login(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToLoginResponse(token), nil
}

func (h *CourierAuthServiceHandler) Authenticate(
	ctx context.Context,
	req *courierv1.AuthenticateRequest,
) (*courierv1.AuthenticateResponse, error) {
	token := req.GetToken()

	courierID, err := h.usecase.Authenticate(ctx, token)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToAuthenticateResponse(courierID), nil
}

var _ courierv1.CourierAuthServiceServer = (*CourierAuthServiceHandler)(nil)
