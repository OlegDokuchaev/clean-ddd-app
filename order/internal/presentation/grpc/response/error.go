package response

import (
	"errors"
	orderDomain "order/internal/domain/order"
	orderRepository "order/internal/infrastructure/repository/order"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParseError(err error) error {
	switch {
	case errors.Is(err, orderDomain.ErrInvalidItems):
	case errors.Is(err, orderDomain.ErrInvalidAddress):
	case errors.Is(err, orderDomain.ErrUnsupportedStatusTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, orderRepository.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, orderRepository.ErrOrderAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	}
	return ErrInternalError
}

var (
	ErrInvalidID     = status.Error(codes.InvalidArgument, "invalid id")
	ErrInternalError = status.Error(codes.Internal, "internal error")
)
