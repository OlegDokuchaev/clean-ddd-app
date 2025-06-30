package customer

import (
	customerGRPC "api-gateway/gen/customer/v1"
	customerDto "api-gateway/internal/domain/dtos/customer"
)

func toRegisterRequest(data customerDto.RegisterDto) *customerGRPC.RegisterRequest {
	return &customerGRPC.RegisterRequest{
		Name:     data.Name,
		Password: data.Password,
		Phone:    data.Phone,
	}
}

func toLoginRequest(data customerDto.LoginDto) *customerGRPC.LoginRequest {
	return &customerGRPC.LoginRequest{
		Phone:    data.Phone,
		Password: data.Password,
	}
}

func toAuthenticateRequest(token string) *customerGRPC.AuthenticateRequest {
	return &customerGRPC.AuthenticateRequest{
		Token: token,
	}
}
