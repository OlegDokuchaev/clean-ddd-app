package usecase

import (
	orderDomain "order/internal/domain/order"

	"github.com/google/uuid"
)

type CreateDto struct {
	CustomerID uuid.UUID
	Address    string
	Items      []orderDomain.Item
}

type BeginDeliveryDto struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
}
