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

func ParseError(err error) error {
	switch {
	case errors.Is(err, courierDomain.ErrInvalidCourierName):
	case errors.Is(err, courierDomain.ErrInvalidCourierPhone):
	case errors.Is(err, courierDomain.ErrInvalidCourierPassword):
	case errors.Is(err, courierRepository.ErrCourierPhoneAlreadyExists):
	case errors.Is(err, courierApplication.ErrAvailableCourierNotFound):
	case errors.Is(err, auth.ErrInvalidSigningMethod):
	case errors.Is(err, auth.ErrInvalidToken):
	case errors.Is(err, auth.ErrTokenExpired):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, courierRepository.ErrCourierNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, courierRepository.ErrCourierAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	}
	return ErrInternalError
}

var (
	ErrInternalError = status.Error(codes.Internal, "internal error")
)
