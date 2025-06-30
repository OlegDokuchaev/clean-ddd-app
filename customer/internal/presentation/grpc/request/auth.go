package request

import (
	customerAuthApplication "customer/internal/application/customer"
	customerv1 "customer/internal/presentation/grpc"
)

func ToRegisterDto(req *customerv1.RegisterRequest) (customerAuthApplication.RegisterDto, error) {
	return customerAuthApplication.RegisterDto{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}, nil
}

func ToLoginDto(req *customerv1.LoginRequest) (customerAuthApplication.LoginDto, error) {
	return customerAuthApplication.LoginDto{
		Phone:    req.Phone,
		Password: req.Password,
	}, nil
}
