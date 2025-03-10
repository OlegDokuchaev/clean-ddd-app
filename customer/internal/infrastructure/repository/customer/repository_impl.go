package customer

import (
	"context"
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/db/tables"

	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, customer *customerDomain.Customer) error {
	model := ToModel(customer)
	res := r.db.WithContext(ctx).Create(model)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByPhone(ctx context.Context, phone string) (*customerDomain.Customer, error) {
	var model tables.Customer
	res := r.db.WithContext(ctx).First(&model, "phone = ?", phone)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

var _ customerDomain.Repository = (*RepositoryImpl)(nil)
