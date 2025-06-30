package item

import (
	"context"
	itemDomain "warehouse/internal/domain/item"
	"warehouse/internal/infrastructure/db/tables"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, item *itemDomain.Item) error {
	model := ToModel(item)
	res := r.db.WithContext(ctx).Omit("Product").Create(model)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) Update(ctx context.Context, item *itemDomain.Item) error {
	res := r.db.WithContext(ctx).Model(&tables.Item{}).
		Where("id = ? AND version = ?", item.ID, item.Version).
		Updates(map[string]any{
			"count":   item.Count,
			"version": uuid.New(),
		})
	if res.RowsAffected == 0 {
		return ErrItemNotFound
	}
	if res.Error != nil {
		item.Version = uuid.New()
	}
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByID(ctx context.Context, itemID uuid.UUID) (*itemDomain.Item, error) {
	var model tables.Item
	res := r.db.WithContext(ctx).Preload("Product").First(&model, "id = ?", itemID)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

func (r *RepositoryImpl) GetAll(ctx context.Context) ([]*itemDomain.Item, error) {
	var models []*tables.Item
	res := r.db.WithContext(ctx).Preload("Product").Find(&models)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomains(models), nil
}

func (r *RepositoryImpl) GetAllByProductIDs(ctx context.Context, productIDs ...uuid.UUID) ([]*itemDomain.Item, error) {
	var models []*tables.Item
	res := r.db.WithContext(ctx).Preload("Product").Find(&models, "product_id IN (?)", productIDs)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	if res.RowsAffected != int64(len(productIDs)) {
		return nil, ErrItemsNotFound
	}
	return ToDomains(models), nil
}

var _ itemDomain.Repository = (*RepositoryImpl)(nil)
