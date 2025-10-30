package customer_request

import (
	customerDto "api-gateway/internal/domain/dtos/customer"
)

func ToRegisterDto(req *RegisterRequest) customerDto.RegisterDto {
	return customerDto.RegisterDto{
		Name:     req.Name,
		Password: req.Password,
		Phone:    req.Phone,
	}
}

func ToLoginDto(req *LoginRequest) customerDto.LoginDto {
	return customerDto.LoginDto{
		Phone:    req.Phone,
		Password: req.Password,
	}
}

func ToVerifyOtp(challengeID string, req *VerifyOtpRequest) customerDto.VerifyOtpDto {
	return customerDto.VerifyOtpDto{
		ChallengeID: challengeID,
		Code:        req.Code,
	}
}
