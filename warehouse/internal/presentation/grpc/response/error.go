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

var domainErrorMap = []struct {
	target error
	code   codes.Code
}{
	// FailedPrecondition
	{productDomain.ErrInvalidProductPrice, codes.FailedPrecondition},
	{productDomain.ErrInvalidProductName, codes.FailedPrecondition},
	{itemDomain.ErrInvalidItemCount, codes.FailedPrecondition},
	{outboxDomain.ErrInvalidOutboxPayload, codes.FailedPrecondition},
	{outboxPublisher.ErrInvalidOutboxMessage, codes.FailedPrecondition},

	// NotFound
	{productRepository.ErrProductNotFound, codes.NotFound},
	{itemRepository.ErrItemNotFound, codes.NotFound},
	{itemRepository.ErrItemsNotFound, codes.NotFound},
	{outboxRepository.ErrOutboxMessageNotFound, codes.NotFound},

	// AlreadyExists
	{productRepository.ErrProductAlreadyExists, codes.AlreadyExists},
	{itemRepository.ErrItemAlreadyExists, codes.AlreadyExists},
	{outboxRepository.ErrOutboxMessageAlreadyExists, codes.AlreadyExists},
}

func ParseError(err error) error {
	for _, e := range domainErrorMap {
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
