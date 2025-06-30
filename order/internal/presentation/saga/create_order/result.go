package create_order

import (
	"encoding/json"
	"github.com/google/uuid"
)

const (
	ItemsReservedName           ResMessageName = "warehouse.items_reserved"
	ItemsReservationFailedName  ResMessageName = "warehouse.items_reservation_failed"
	ItemsReleasedName           ResMessageName = "warehouse.items_released"
	CourierAssignedName         ResMessageName = "courier.courier_assigned"
	CourierAssignmentFailedName ResMessageName = "courier.courier_assignment_failed"
)

type (
	ResMessageName string

	ResMessage struct {
		ID      uuid.UUID
		Name    ResMessageName
		Payload json.RawMessage
	}
)
