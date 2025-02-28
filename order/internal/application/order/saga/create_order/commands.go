package create_order

import (
	"github.com/google/uuid"
	orderDomain "order/internal/domain/order"
)

type Cmd struct {
	ID uuid.UUID
}

type ReserveItemsCmd struct {
	OrderID uuid.UUID
	Items   []orderDomain.Item
	Cmd
}

type ReleaseItemsCmd struct {
	OrderID uuid.UUID
	Items   []orderDomain.Item
	Cmd
}

type CancelOutOfStockCmd struct {
	OrderID uuid.UUID
	Cmd
}

type AssignCourierCmd struct {
	OrderID uuid.UUID
	Cmd
}

type BeginDeliveryCmd struct {
	OrderID uuid.UUID
	Cmd
}

type CancelCourierNotFoundCmd struct {
	OrderID uuid.UUID
	Cmd
}
