package customer

import (
	"context"
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/db/tables"
	"github.com/google/uuid"

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

func (r *RepositoryImpl) Save(ctx context.Context, customer *customerDomain.Customer) error {
	model := ToModel(customer)
	res := r.db.WithContext(ctx).Save(model)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*customerDomain.Customer, error) {
	var model tables.Customer
	res := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

func (r *RepositoryImpl) GetByPhone(ctx context.Context, phone string) (*customerDomain.Customer, error) {
	var model tables.Customer
	res := r.db.WithContext(ctx).First(&model, "phone = ?", phone)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

func (r *RepositoryImpl) GetByEmail(ctx context.Context, email string) (*customerDomain.Customer, error) {
	var model tables.Customer
	res := r.db.WithContext(ctx).First(&model, "email = ?", email)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

var _ customerDomain.Repository = (*RepositoryImpl)(nil)
