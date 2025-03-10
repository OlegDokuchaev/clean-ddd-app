package order

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	orderDomain "order/internal/domain/order"
	"order/internal/infrastructure/db/tables"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, order *orderDomain.Order) error {
	orderModel := ToModel(order)
	res := r.db.WithContext(ctx).Create(orderModel)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) Update(ctx context.Context, order *orderDomain.Order) error {
	newVersion := uuid.New()
	orderModel := ToModel(order)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&tables.Order{}).
			Where("id = ? AND version = ?", orderModel.ID, orderModel.Version).
			Updates(map[string]any{
				"customer_id": orderModel.CustomerID,
				"status":      orderModel.Status,
				"created":     orderModel.Created,
				"version":     newVersion,
			})
		if res.Error != nil {
			return ParseError(res.Error)
		}
		if res.RowsAffected == 0 {
			return ErrOrderNotFound
		}

		return updateDelivery(tx, orderModel.ID, orderModel.Delivery)
	})

	if err == nil {
		order.Version = newVersion
	}
	return err
}

func (r *RepositoryImpl) GetByID(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
	var orderModel tables.Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Delivery").
		Where("id = ?", orderID).
		First(&orderModel).Error; err != nil {
		return nil, ParseError(err)
	}
	return ToDomain(&orderModel), nil
}

func (r *RepositoryImpl) GetAllByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderDomain.Order, error) {
	var orderModels []*tables.Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Delivery").
		Where("customer_id = ?", customerID).
		Find(&orderModels).Error; err != nil {
		return nil, ParseError(err)
	}
	return ToDomains(orderModels), nil
}

func (r *RepositoryImpl) GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderDomain.Order, error) {
	var orderModels []*tables.Order
	if err := r.db.WithContext(ctx).
		Joins("JOIN deliveries ON deliveries.order_id = orders.id").
		Preload("Items").
		Preload("Delivery").
		Where("deliveries.courier_id = ?", courierID).
		Find(&orderModels).Error; err != nil {
		return nil, ParseError(err)
	}
	return ToDomains(orderModels), nil
}

func updateDelivery(tx *gorm.DB, orderID uuid.UUID, delivery tables.Delivery) error {
	if err := tx.Where("order_id = ?", orderID).Delete(&tables.Delivery{}).Error; err != nil {
		return ParseError(err)
	}
	res := tx.Create(&delivery)
	return ParseError(res.Error)
}

var _ orderDomain.Repository = (*RepositoryImpl)(nil)
