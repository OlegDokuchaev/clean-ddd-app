package commands

import (
	"github.com/google/uuid"
)

const (
	ItemsReservedName          ResMessageName = "warehouse.items_reserved"
	ItemsReservationFailedName ResMessageName = "warehouse.items_reservation_failed"
	ItemsReleasedName          ResMessageName = "warehouse.items_released"
)

type (
	ResMessageName    string
	ResMessagePayload interface{}

	ResMessage struct {
		ID      uuid.UUID
		Name    ResMessageName
		Payload ResMessagePayload
	}
)

type ItemsReserved struct {
	OrderID uuid.UUID
}

type ItemsReservationFailed struct {
	OrderID uuid.UUID
}

type ItemsReleased struct {
	OrderID uuid.UUID
}

func newResMessage(name ResMessageName, payload ResMessagePayload) *ResMessage {
	return &ResMessage{
		ID:      uuid.New(),
		Name:    name,
		Payload: payload,
	}
}
