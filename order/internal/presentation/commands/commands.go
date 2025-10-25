package commands

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	CancelOutOfStockCmdName      CmdMessageName = "create_order.cancel_out_of_stock"
	BeginDeliveryCmdName         CmdMessageName = "create_order.begin_delivery"
	CancelCourierNotFoundCmdName CmdMessageName = "create_order.cancel_courier_not_found"
)

type (
	CmdMessageName string

	CmdMessage struct {
		ID      uuid.UUID
		Name    CmdMessageName
		Payload json.RawMessage
	}
)

type CancelOutOfStockCmd struct {
	OrderID uuid.UUID
}

type CancelCourierNotFoundCmd struct {
	OrderID uuid.UUID
}

type BeginDeliveryCmd struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
}
