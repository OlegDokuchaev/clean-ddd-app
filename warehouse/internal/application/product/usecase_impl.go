package product

import (
	"context"
	"github.com/google/uuid"
	"warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/domain/uow"
)

type UseCaseImpl struct {
	uow uow.UoW
}

func (u *UseCaseImpl) Create(ctx context.Context, data CreateDto) (uuid.UUID, error) {
	product, events, err := productDomain.Create(data.Name, data.Price)
	if err != nil {
		return uuid.Nil, err
	}

	messages := make([]*outbox.Message, 0, len(events))
	for _, event := range events {
		message, err := outbox.Create(event)
		if err != nil {
			return uuid.Nil, err
		}
		messages = append(messages, message)
	}

	err = u.uow.Transaction(ctx, func(u uow.UoW) error {
		if err = u.Product().Create(ctx, product); err != nil {
			return err
		}

		for _, message := range messages {
			if err = u.Outbox().Create(ctx, message); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
