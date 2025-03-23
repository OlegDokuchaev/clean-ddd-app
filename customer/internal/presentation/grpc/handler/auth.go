package handler

import (
	"context"
	customerApplication "customer/internal/application/customer"
	customerv1 "customer/internal/presentation/grpc"
	"customer/internal/presentation/grpc/request"
	"customer/internal/presentation/grpc/response"
)

type CustomerAuthServiceHandler struct {
	customerv1.UnimplementedCustomerAuthServiceServer

	usecase customerApplication.AuthUseCase
}

func NewCustomerAuthServiceHandler(usecase customerApplication.AuthUseCase) *CustomerAuthServiceHandler {
	return &CustomerAuthServiceHandler{
		usecase: usecase,
	}
}

func (h *CustomerAuthServiceHandler) Register(
	ctx context.Context,
	req *customerv1.RegisterRequest,
) (*customerv1.RegisterResponse, error) {
	data, err := request.ToRegisterDto(req)
	if err != nil {
		return nil, err
	}

	customerID, err := h.usecase.Register(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToRegisterResponse(customerID), nil
}

func (h *CustomerAuthServiceHandler) Login(
	ctx context.Context,
	req *customerv1.LoginRequest,
) (*customerv1.LoginResponse, error) {
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

func (h *CustomerAuthServiceHandler) Authenticate(
	ctx context.Context,
	req *customerv1.AuthenticateRequest,
) (*customerv1.AuthenticateResponse, error) {
	token := req.GetToken()

	customerID, err := h.usecase.Authenticate(ctx, token)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToAuthenticateResponse(customerID), nil
}

var _ customerv1.CustomerAuthServiceServer = (*CustomerAuthServiceHandler)(nil)
