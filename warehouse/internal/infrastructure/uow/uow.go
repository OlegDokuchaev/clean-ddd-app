package uow

import (
	"context"
	itemDomain "warehouse/internal/domain/item"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/domain/uow"
	itemRepository "warehouse/internal/infrastructure/repository/item"
	outboxRepository "warehouse/internal/infrastructure/repository/outbox"
	productRepository "warehouse/internal/infrastructure/repository/product"

	"gorm.io/gorm"
)

type UoWImpl struct {
	productRepository productDomain.Repository
	itemRepository    itemDomain.Repository
	outboxRepository  outboxDomain.Repository

	db *gorm.DB
}

func New(db *gorm.DB) uow.UoW {
	return &UoWImpl{
		productRepository: productRepository.New(db),
		itemRepository:    itemRepository.New(db),
		outboxRepository:  outboxRepository.New(db),
		db:                db,
	}
}

func (u *UoWImpl) Transaction(ctx context.Context, fn func(uow.UoW) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txManager := New(tx)
		return fn(txManager)
	})
}

func (u *UoWImpl) Product() productDomain.Repository {
	return u.productRepository
}

func (u *UoWImpl) Item() itemDomain.Repository {
	return u.itemRepository
}

func (u *UoWImpl) Outbox() outboxDomain.Repository {
	return u.outboxRepository
}

var _ uow.UoW = (*UoWImpl)(nil)
