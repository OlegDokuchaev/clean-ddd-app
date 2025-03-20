package create_order

import (
	"github.com/google/uuid"
)

type ReserveItemsCmd struct {
	OrderID uuid.UUID
	Items   []OrderItem
}

type ReleaseItemsCmd struct {
	OrderID uuid.UUID
	Items   []OrderItem
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

type OrderItem struct {
	ProductID uuid.UUID
	Count     int
}
