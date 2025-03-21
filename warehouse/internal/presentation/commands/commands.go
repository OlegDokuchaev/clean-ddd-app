package commands

import (
	itemApplication "warehouse/internal/application/item"

	"github.com/google/uuid"
)

const (
	ReserveItemsCmdName CmdMessageName = "create_order.reserve_items"
	ReleaseItemsCmdName CmdMessageName = "create_order.release_items"
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

type ReserveItemsCmd struct {
	OrderID uuid.UUID
	Items   []itemApplication.ItemDto
}

type ReleaseItemsCmd struct {
	OrderID uuid.UUID
	Items   []itemApplication.ItemDto
}
