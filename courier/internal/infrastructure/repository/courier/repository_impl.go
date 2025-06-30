package courier

import (
	"context"
	courierDomain "courier/internal/domain/courier"
	"courier/internal/infrastructure/db/tables"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, courier *courierDomain.Courier) error {
	model := ToModel(courier)
	res := r.db.WithContext(ctx).Create(model)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByPhone(ctx context.Context, phone string) (*courierDomain.Courier, error) {
	var model tables.Courier
	res := r.db.WithContext(ctx).First(&model, "phone = ?", phone)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

func (r *RepositoryImpl) GetAll(ctx context.Context) ([]*courierDomain.Courier, error) {
	var models []*tables.Courier
	res := r.db.WithContext(ctx).Find(&models)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomains(models), nil
}

var _ courierDomain.Repository = (*RepositoryImpl)(nil)
