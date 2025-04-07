package request

import (
	courierDto "api-gateway/internal/domain/dtos/courier"
)

func ToRegisterDto(req *RegisterRequest) courierDto.RegisterDto {
	return courierDto.RegisterDto{
		Name:     req.Name,
		Password: req.Password,
		Phone:    req.Phone,
	}
}

func ToLoginDto(req *LoginRequest) courierDto.LoginDto {
	return courierDto.LoginDto{
		Phone:    req.Phone,
		Password: req.Password,
	}
}
