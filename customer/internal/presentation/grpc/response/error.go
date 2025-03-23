package response

import (
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/auth"
	customerRepository "customer/internal/infrastructure/repository/customer"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParseError(err error) error {
	switch {
	case errors.Is(err, customerDomain.ErrInvalidCustomerName):
	case errors.Is(err, customerDomain.ErrInvalidCustomerPhone):
	case errors.Is(err, customerDomain.ErrInvalidCustomerPassword):
	case errors.Is(err, customerRepository.ErrCustomerPhoneAlreadyExists):
	case errors.Is(err, auth.ErrInvalidSigningMethod):
	case errors.Is(err, auth.ErrInvalidToken):
	case errors.Is(err, auth.ErrTokenExpired):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, customerRepository.ErrCustomerNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, customerRepository.ErrCustomerAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	}
	return ErrInternalError
}

var (
	ErrInternalError = status.Error(codes.Internal, "internal error")
)
