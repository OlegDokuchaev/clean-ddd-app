package response

import (
	"errors"
	orderDomain "order/internal/domain/order"
	orderRepository "order/internal/infrastructure/repository/order"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var orderErrorMap = []struct {
	target error
	code   codes.Code
}{
	// InvalidArgument
	{orderDomain.ErrInvalidItems, codes.InvalidArgument},
	{orderDomain.ErrInvalidAddress, codes.InvalidArgument},
	{orderDomain.ErrUnsupportedStatusTransition, codes.InvalidArgument},

	// NotFound
	{orderRepository.ErrOrderNotFound, codes.NotFound},

	// AlreadyExists
	{orderRepository.ErrOrderAlreadyExists, codes.AlreadyExists},
}

func ParseError(err error) error {
	for _, e := range orderErrorMap {
		if errors.Is(err, e.target) {
			return status.Error(e.code, err.Error())
		}
	}
	return ErrInternalError
}

var (
	ErrInvalidID     = status.Error(codes.InvalidArgument, "invalid id")
	ErrInternalError = status.Error(codes.Internal, "internal error")
)
