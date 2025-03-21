package item

import (
	"context"
	itemDomain "warehouse/internal/domain/item"
	"warehouse/internal/domain/uow"
	itemRepository "warehouse/internal/infrastructure/repository/item"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	uow uow.UoW
}

func NewUseCase(uow uow.UoW) UseCase {
	return &UseCaseImpl{uow: uow}
}

func (u *UseCaseImpl) Create(ctx context.Context, data CreateDto) (uuid.UUID, error) {
	product, err := u.uow.Product().GetByID(ctx, data.ProductID)
	if err != nil {
		return uuid.Nil, err
	}

	item, err := itemDomain.Create(product, data.Count)
	if err != nil {
		return uuid.Nil, err
	}

	err = u.uow.Item().Create(ctx, item)
	if err != nil {
		return uuid.Nil, err
	}

	return item.ID, nil
}

func (u *UseCaseImpl) Reserve(ctx context.Context, data ReserveDto) error {
	productIDs := make([]uuid.UUID, 0, len(data.Items))
	for _, itemDto := range data.Items {
		productIDs = append(productIDs, itemDto.ProductID)
	}

	items, err := u.uow.Item().GetAllByProductIDs(ctx, productIDs...)
	if err != nil {
		return err
	}

	reserveMap := make(map[uuid.UUID]int, len(data.Items))
	for _, itemDto := range data.Items {
		reserveMap[itemDto.ProductID] = itemDto.Count
	}

	return u.uow.Transaction(ctx, func(tx uow.UoW) error {
		for _, item := range items {
			count, exists := reserveMap[item.Product.ID]
			if !exists {
				return itemRepository.ErrItemsNotFound
			}

			if err := item.Reserve(count); err != nil {
				return err
			}

			if err := tx.Item().Update(ctx, item); err != nil {
				return err
			}
		}

		return nil
	})
}

func (u *UseCaseImpl) Release(ctx context.Context, data ReleaseDto) error {
	productIDs := make([]uuid.UUID, 0, len(data.Items))
	for _, itemDto := range data.Items {
		productIDs = append(productIDs, itemDto.ProductID)
	}

	items, err := u.uow.Item().GetAllByProductIDs(ctx, productIDs...)
	if err != nil {
		return err
	}

	releaseMap := make(map[uuid.UUID]int, len(data.Items))
	for _, itemDto := range data.Items {
		releaseMap[itemDto.ProductID] = itemDto.Count
	}

	return u.uow.Transaction(ctx, func(tx uow.UoW) error {
		for _, item := range items {
			count, exists := releaseMap[item.Product.ID]
			if !exists {
				return itemRepository.ErrItemsNotFound
			}

			if err := item.Release(count); err != nil {
				return err
			}

			if err := tx.Item().Update(ctx, item); err != nil {
				return err
			}
		}

		return nil
	})
}

func (u *UseCaseImpl) GetAll(ctx context.Context) ([]*itemDomain.Item, error) {
	return u.uow.Item().GetAll(ctx)
}

var _ UseCase = (*UseCaseImpl)(nil)
