package response

import (
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/auth"
	customerRepository "customer/internal/infrastructure/repository/customer"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var customerErrorMap = []struct {
	target error
	code   codes.Code
}{
	// InvalidArgument
	{customerDomain.ErrInvalidCustomerName, codes.InvalidArgument},
	{customerDomain.ErrInvalidCustomerPhone, codes.InvalidArgument},
	{customerDomain.ErrInvalidCustomerPassword, codes.InvalidArgument},
	{customerRepository.ErrCustomerPhoneAlreadyExists, codes.InvalidArgument},
	{auth.ErrInvalidSigningMethod, codes.InvalidArgument},
	{auth.ErrInvalidToken, codes.InvalidArgument},
	{auth.ErrTokenExpired, codes.InvalidArgument},

	// NotFound
	{customerRepository.ErrCustomerNotFound, codes.NotFound},

	// AlreadyExists
	{customerRepository.ErrCustomerAlreadyExists, codes.AlreadyExists},
}

func ParseError(err error) error {
	for _, e := range customerErrorMap {
		if errors.Is(err, e.target) {
			return status.Error(e.code, err.Error())
		}
	}
	return ErrInternalError
}

var ErrInternalError = status.Error(codes.Internal, "internal error")
