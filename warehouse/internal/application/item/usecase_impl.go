package item

import (
	"context"
	domainItem "warehouse/internal/domain/item"
	"warehouse/internal/domain/uow"

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

	item, err := domainItem.Create(product, data.Count)
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
	item, err := u.uow.Item().GetByID(ctx, data.ItemID)
	if err != nil {
		return err
	}

	if err = item.Reserve(data.Count); err != nil {
		return err
	}
	if err = u.uow.Item().Update(ctx, item); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) Release(ctx context.Context, data ReleaseDto) error {
	item, err := u.uow.Item().GetByID(ctx, data.ItemID)
	if err != nil {
		return err
	}

	if err = item.Release(data.Count); err != nil {
		return err
	}
	if err = u.uow.Item().Update(ctx, item); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) GetAll(ctx context.Context) ([]*domainItem.Item, error) {
	return u.uow.Item().GetAll(ctx)
}

var _ UseCase = (*UseCaseImpl)(nil)
