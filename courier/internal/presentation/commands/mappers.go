package commands

import "github.com/google/uuid"

func toCourierAssignmentFailed(orderID uuid.UUID) *ResMessage {
	return newResMessage(CourierAssignmentFailedName, CourierAssignmentFailed{
		OrderID: orderID,
	})
}

func toCourierAssigned(orderID uuid.UUID, courierID uuid.UUID) *ResMessage {
	return newResMessage(CourierAssignedName, CourierAssigned{
		OrderID:   orderID,
		CourierID: courierID,
	})
}
