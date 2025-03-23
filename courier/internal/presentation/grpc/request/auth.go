package request

import (
	courierAuthApplication "courier/internal/application/courier/auth"
	courierv1 "courier/internal/presentation/grpc"
)

func ToRegisterDto(req *courierv1.RegisterRequest) (courierAuthApplication.RegisterDto, error) {
	return courierAuthApplication.RegisterDto{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}, nil
}

func ToLoginDto(req *courierv1.LoginRequest) (courierAuthApplication.LoginDto, error) {
	return courierAuthApplication.LoginDto{
		Phone:    req.Phone,
		Password: req.Password,
	}, nil
}
