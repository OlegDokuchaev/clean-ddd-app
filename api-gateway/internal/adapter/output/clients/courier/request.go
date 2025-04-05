package courier

import (
	courierGRPC "api-gateway/gen/courier/v1"
	courierDto "api-gateway/internal/domain/dtos/courier"
)

func toRegisterRequest(data courierDto.RegisterDto) *courierGRPC.RegisterRequest {
	return &courierGRPC.RegisterRequest{
		Name:     data.Name,
		Password: data.Password,
		Phone:    data.Phone,
	}
}

func toLoginRequest(data courierDto.LoginDto) *courierGRPC.LoginRequest {
	return &courierGRPC.LoginRequest{
		Phone:    data.Phone,
		Password: data.Password,
	}
}

func toAuthenticateRequest(token string) *courierGRPC.AuthenticateRequest {
	return &courierGRPC.AuthenticateRequest{
		Token: token,
	}
}
