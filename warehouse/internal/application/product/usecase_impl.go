package product

import (
	"context"
	"warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/domain/uow"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	uow          uow.UoW
	imageService ImageService
}

func NewUseCase(uow uow.UoW, imageService ImageService) *UseCaseImpl {
	return &UseCaseImpl{
		uow:          uow,
		imageService: imageService,
	}
}

func (u *UseCaseImpl) Create(ctx context.Context, data CreateDto) (uuid.UUID, error) {
	path, err := u.imageService.Create(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	product, events, err := productDomain.Create(data.Name, data.Price, path)
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
