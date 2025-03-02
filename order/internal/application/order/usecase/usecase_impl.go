package usecase

import (
	"context"
	"github.com/google/uuid"
	createOrderSaga "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
)

type UseCaseImpl struct {
	repo                   orderDomain.Repository
	createOrderSagaManager createOrderSaga.Manager
}

func New(repo orderDomain.Repository, createOrderSagaManager createOrderSaga.Manager) UseCase {
	return &UseCaseImpl{
		repo:                   repo,
		createOrderSagaManager: createOrderSagaManager,
	}
}

func (u *UseCaseImpl) Create(ctx context.Context, data CreateDto) (uuid.UUID, error) {
	order := orderDomain.Create(data.CustomerID, data.Address, data.Items)

	if err := u.repo.Create(ctx, order); err != nil {
		return uuid.Nil, err
	}
	u.createOrderSagaManager.Create(ctx, order)

	return order.ID, nil
}

func (u *UseCaseImpl) CancelByCustomer(ctx context.Context, orderID uuid.UUID) error {
	order, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err = order.NoteCanceledByCustomer(); err != nil {
		return err
	}
	if err = u.repo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) CancelOutOfStock(ctx context.Context, orderID uuid.UUID) error {
	order, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err = order.NoteCanceledOutOfStock(); err != nil {
		return err
	}
	if err = u.repo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) CancelCourierNotFound(ctx context.Context, orderID uuid.UUID) error {
	order, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err = order.NoteCanceledCourierNotFound(); err != nil {
		return err
	}
	if err = u.repo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) BeginDelivery(ctx context.Context, data BeginDeliveryDto) error {
	order, err := u.repo.GetByID(ctx, data.OrderID)
	if err != nil {
		return err
	}

	if err = order.NoteDelivering(data.CourierID); err != nil {
		return err
	}
	if err = u.repo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) CompleteDelivery(ctx context.Context, orderID uuid.UUID) error {
	order, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err = order.NoteDelivered(); err != nil {
		return err
	}
	if err = u.repo.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) GetAllByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderDomain.Order, error) {
	return u.repo.GetAllByCustomer(ctx, customerID)
}

func (u *UseCaseImpl) GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderDomain.Order, error) {
	return u.repo.GetCurrentByCourier(ctx, courierID)
}
