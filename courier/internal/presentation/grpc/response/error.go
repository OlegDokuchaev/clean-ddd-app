package response

import (
	courierApplication "courier/internal/application/courier"
	courierDomain "courier/internal/domain/courier"
	"courier/internal/infrastructure/auth"
	courierRepository "courier/internal/infrastructure/repository/courier"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var parseErrorMap = []struct {
	target error
	code   codes.Code
}{
	// InvalidArgument
	{courierDomain.ErrInvalidCourierName, codes.InvalidArgument},
	{courierDomain.ErrInvalidCourierPhone, codes.InvalidArgument},
	{courierDomain.ErrInvalidCourierPassword, codes.InvalidArgument},
	{courierRepository.ErrCourierPhoneAlreadyExists, codes.InvalidArgument},
	{courierApplication.ErrAvailableCourierNotFound, codes.InvalidArgument},
	{auth.ErrInvalidSigningMethod, codes.InvalidArgument},
	{auth.ErrInvalidToken, codes.InvalidArgument},
	{auth.ErrTokenExpired, codes.InvalidArgument},

	// NotFound
	{courierRepository.ErrCourierNotFound, codes.NotFound},

	// AlreadyExists
	{courierRepository.ErrCourierAlreadyExists, codes.AlreadyExists},
}

func ParseError(err error) error {
	for _, e := range parseErrorMap {
		if errors.Is(err, e.target) {
			return status.Error(e.code, err.Error())
		}
	}
	return ErrInternalError
}

var ErrInternalError = status.Error(codes.Internal, "internal error")
