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
	{courierDomain.ErrInvalidCourierName, codes.FailedPrecondition},
	{courierDomain.ErrInvalidCourierPhone, codes.FailedPrecondition},
	{courierDomain.ErrInvalidCourierPassword, codes.FailedPrecondition},
	{courierRepository.ErrCourierPhoneAlreadyExists, codes.FailedPrecondition},
	{courierApplication.ErrAvailableCourierNotFound, codes.FailedPrecondition},
	{auth.ErrInvalidSigningMethod, codes.FailedPrecondition},
	{auth.ErrInvalidToken, codes.FailedPrecondition},
	{auth.ErrTokenExpired, codes.FailedPrecondition},

	{courierRepository.ErrCourierNotFound, codes.NotFound},

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
