package order

import (
	orderDto "api-gateway/internal/domain/dtos/order"
	courierClient "api-gateway/internal/port/output/clients/courier"
	customerClient "api-gateway/internal/port/output/clients/customer"
	orderClient "api-gateway/internal/port/output/clients/order"
	"context"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	customerClient customerClient.Client
	courierClient  courierClient.Client
	orderClient    orderClient.Client
}

func NewUseCase(
	customerClient customerClient.Client,
	courierClient courierClient.Client,
	orderClient orderClient.Client,
) UseCase {
	return &UseCaseImpl{
		customerClient: customerClient,
		courierClient:  courierClient,
		orderClient:    orderClient,
	}
}

func (u *UseCaseImpl) Create(ctx context.Context, data orderDto.CreateDto, customerToken string) (uuid.UUID, error) {
	customerID, err := u.customerClient.Authenticate(ctx, customerToken)
	if err != nil {
		return uuid.Nil, err
	}

	dto := orderClient.CreateDto{
		CustomerID: customerID,
		Address:    data.Address,
		Items:      data.Items,
	}
	orderID, err := u.orderClient.Create(ctx, dto)
	if err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

func (u *UseCaseImpl) CancelByCustomer(ctx context.Context, orderID uuid.UUID, customerToken string) error {
	customerID, err := u.customerClient.Authenticate(ctx, customerToken)
	if err != nil {
		return err
	}

	err = u.orderClient.CancelByCustomer(ctx, orderID, customerID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) Complete(ctx context.Context, orderID uuid.UUID, courierToken string) error {
	courierID, err := u.courierClient.Authenticate(ctx, courierToken)
	if err != nil {
		return err
	}

	err = u.orderClient.Complete(ctx, orderID, courierID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) GetByCustomer(ctx context.Context, customerToken string) ([]*orderDto.OrderDto, error) {
	customerID, err := u.customerClient.Authenticate(ctx, customerToken)
	if err != nil {
		return nil, err
	}

	orders, err := u.orderClient.GetByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (u *UseCaseImpl) GetCurrentByCourier(ctx context.Context, courierToken string) ([]*orderDto.OrderDto, error) {
	courierID, err := u.courierClient.Authenticate(ctx, courierToken)
	if err != nil {
		return nil, err
	}

	orders, err := u.orderClient.GetCurrentByCourier(ctx, courierID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
