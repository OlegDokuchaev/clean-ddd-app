package create_order

import (
	"github.com/google/uuid"
	orderDomain "order/internal/domain/order"
)

type ReserveItemsCmd struct {
	OrderID uuid.UUID
	Items   []orderDomain.Item
}

type ReleaseItemsCmd struct {
	OrderID uuid.UUID
	Items   []orderDomain.Item
}

type CancelOutOfStockCmd struct {
	OrderID uuid.UUID
}

type AssignCourierCmd struct {
	OrderID uuid.UUID
}

type BeginDeliveryCmd struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
}

type CancelCourierNotFoundCmd struct {
	OrderID uuid.UUID
}
