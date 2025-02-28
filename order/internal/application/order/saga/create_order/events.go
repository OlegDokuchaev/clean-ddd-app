package create_order

import "github.com/google/uuid"

type Event struct {
	ID uuid.UUID
}

type ItemsReserved struct {
	OrderID uuid.UUID
	Event
}

type ItemsReservationFailed struct {
	OrderID uuid.UUID
	Event
}

type CourierAssignmentFailed struct {
	OrderID uuid.UUID
	Event
}

type CourierAssigned struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
	Event
}

type ItemsReleased struct {
	OrderID uuid.UUID
	Event
}
