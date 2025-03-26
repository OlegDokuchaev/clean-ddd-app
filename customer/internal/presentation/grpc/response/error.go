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
	{customerDomain.ErrInvalidCustomerName, codes.FailedPrecondition},
	{customerDomain.ErrInvalidCustomerPhone, codes.FailedPrecondition},
	{customerDomain.ErrInvalidCustomerPassword, codes.FailedPrecondition},
	{customerRepository.ErrCustomerPhoneAlreadyExists, codes.FailedPrecondition},
	{auth.ErrInvalidSigningMethod, codes.FailedPrecondition},
	{auth.ErrInvalidToken, codes.FailedPrecondition},
	{auth.ErrTokenExpired, codes.FailedPrecondition},

	{customerRepository.ErrCustomerNotFound, codes.NotFound},

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
