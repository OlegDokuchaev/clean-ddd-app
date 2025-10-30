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

func toVerifyOtpRequest(data customerDto.VerifyOtpDto) *customerGRPC.VerifyOtpRequest {
	return &customerGRPC.VerifyOtpRequest{
		ChallengeId: data.ChallengeID,
		Code:        data.Code,
	}
}

func toRequestPasswordResetRequest(email string) *customerGRPC.RequestPasswordResetRequest {
	return &customerGRPC.RequestPasswordResetRequest{
		Email: email,
	}
}

func toCompletePasswordResetRequest(token string, newPassword string) *customerGRPC.CompletePasswordResetRequest {
	return &customerGRPC.CompletePasswordResetRequest{
		Token:       token,
		NewPassword: newPassword,
	}
}

func toAuthenticateRequest(token string) *customerGRPC.AuthenticateRequest {
	return &customerGRPC.AuthenticateRequest{
		Token: token,
	}
}
