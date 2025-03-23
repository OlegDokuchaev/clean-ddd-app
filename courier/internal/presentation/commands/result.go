package commands

import (
	"github.com/google/uuid"
)

const (
	CourierAssignmentFailedName ResMessageName = "courier.courier_assignment_failed"
	CourierAssignedName         ResMessageName = "courier.courier_assigned"
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

type CourierAssignmentFailed struct {
	OrderID uuid.UUID
}

type CourierAssigned struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
}

func newResMessage(name ResMessageName, payload ResMessagePayload) *ResMessage {
	return &ResMessage{
		ID:      uuid.New(),
		Name:    name,
		Payload: payload,
	}
}
