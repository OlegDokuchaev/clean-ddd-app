package commands

import (
	itemApplication "warehouse/internal/application/item"

	"github.com/google/uuid"
)

func toReserveItemsDto(cmd ReserveItemsCmd) itemApplication.ReserveDto {
	return itemApplication.ReserveDto{
		Items: cmd.Items,
	}
}

func toReleaseItemsDto(cmd ReleaseItemsCmd) itemApplication.ReleaseDto {
	return itemApplication.ReleaseDto{
		Items: cmd.Items,
	}
}

func toItemsReserved(orderID uuid.UUID) *ResMessage {
	return newResMessage(ItemsReservedName, ItemsReserved{
		OrderID: orderID,
	})
}

func toItemsReservationFailed(orderID uuid.UUID) *ResMessage {
	return newResMessage(ItemsReservationFailedName, ItemsReservationFailed{
		OrderID: orderID,
	})
}

func toItemsReleased(orderID uuid.UUID) *ResMessage {
	return newResMessage(ItemsReleasedName, ItemsReleased{
		OrderID: orderID,
	})
}
