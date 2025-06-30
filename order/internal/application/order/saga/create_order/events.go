package create_order

import "github.com/google/uuid"

type ItemsReserved struct {
	OrderID uuid.UUID
}

type ItemsReservationFailed struct {
	OrderID uuid.UUID
}

type CourierAssignmentFailed struct {
	OrderID uuid.UUID
}

type CourierAssigned struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
}

type ItemsReleased struct {
	OrderID uuid.UUID
}
