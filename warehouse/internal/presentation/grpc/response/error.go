package response

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	itemDomain "warehouse/internal/domain/item"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	outboxPublisher "warehouse/internal/infrastructure/publisher/outbox"
	itemRepository "warehouse/internal/infrastructure/repository/item"
	outboxRepository "warehouse/internal/infrastructure/repository/outbox"
	productRepository "warehouse/internal/infrastructure/repository/product"
)

func ParseError(err error) error {
	switch {
	case errors.Is(err, productDomain.ErrInvalidProductPrice):
	case errors.Is(err, productDomain.ErrInvalidProductName):
	case errors.Is(err, itemDomain.ErrInvalidItemCount):
	case errors.Is(err, outboxDomain.ErrInvalidOutboxPayload):
	case errors.Is(err, outboxPublisher.ErrInvalidOutboxMessage):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, productRepository.ErrProductNotFound):
	case errors.Is(err, itemRepository.ErrItemNotFound):
	case errors.Is(err, itemRepository.ErrItemsNotFound):
	case errors.Is(err, outboxRepository.ErrOutboxMessageNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, productRepository.ErrProductAlreadyExists):
	case errors.Is(err, itemRepository.ErrItemAlreadyExists):
	case errors.Is(err, outboxRepository.ErrOutboxMessageAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	}
	return ErrInternalError
}

var (
	ErrInvalidID     = status.Error(codes.InvalidArgument, "invalid id")
	ErrInternalError = status.Error(codes.Internal, "internal error")
)
