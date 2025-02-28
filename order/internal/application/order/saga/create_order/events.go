package create_order

import "github.com/google/uuid"

type EventBase struct {
	ID uuid.UUID
}

type ItemsReserved struct {
	OrderID uuid.UUID
	EventBase
}

type ItemsReservationFailed struct {
	OrderID uuid.UUID
	EventBase
}

type CourierAssignmentFailed struct {
	OrderID uuid.UUID
	EventBase
}

type CourierAssigned struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
	EventBase
}

type ItemsReleased struct {
	OrderID uuid.UUID
	EventBase
}
