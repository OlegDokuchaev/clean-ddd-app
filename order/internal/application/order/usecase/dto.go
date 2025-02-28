package usecase

import (
	"github.com/google/uuid"
	orderDomain "order/internal/domain/order"
)

type CreateDto struct {
	CustomerID uuid.UUID
	Address    string
	Items      []orderDomain.Item
}

type BeginDeliveryDto struct {
	courierID uuid.UUID
}
