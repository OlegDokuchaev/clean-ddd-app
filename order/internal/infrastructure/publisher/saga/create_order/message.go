package create_order

import "github.com/google/uuid"

const (
	ReserveItemsCmdName          CmdMessageName = "create_order.reserve_items"
	ReleaseItemsCmdName          CmdMessageName = "create_order.release_items"
	CancelOutOfStockCmdName      CmdMessageName = "create_order.cancel_out_of_stock"
	AssignCourierCmdName         CmdMessageName = "create_order.assign_courier"
	BeginDeliveryCmdName         CmdMessageName = "create_order.begin_delivery"
	CancelCourierNotFoundCmdName CmdMessageName = "create_order.cancel_courier_not_found"
)

type (
	CmdMessageName    string
	CmdMessagePayload interface{}

	CmdMessage struct {
		ID      uuid.UUID
		Name    CmdMessageName
		Payload CmdMessagePayload
	}
)

func NewCmdMessage(name CmdMessageName, payload CmdMessagePayload) CmdMessage {
	return CmdMessage{
		ID:      uuid.New(),
		Name:    name,
		Payload: payload,
	}
}
