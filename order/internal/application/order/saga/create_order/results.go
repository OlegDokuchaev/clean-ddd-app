package create_order

import "github.com/google/uuid"

type Result struct {
	ID uuid.UUID
}

type ItemsReserved struct {
	OrderID uuid.UUID
	Result
}

type ItemsReservationFailed struct {
	OrderID uuid.UUID
	Result
}

type CourierAssignmentFailed struct {
	OrderID uuid.UUID
	Result
}

type CourierAssigned struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
	Result
}

type ItemsReleased struct {
	OrderID uuid.UUID
	Result
}
