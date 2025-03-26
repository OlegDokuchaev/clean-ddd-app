package commands

import (
	"encoding/json"
	itemApplication "warehouse/internal/application/item"

	"github.com/google/uuid"
)

const (
	ReserveItemsCmdName CmdMessageName = "create_order.reserve_items"
	ReleaseItemsCmdName CmdMessageName = "create_order.release_items"
)

type (
	CmdMessageName string

	CmdMessage struct {
		ID      uuid.UUID
		Name    CmdMessageName
		Payload json.RawMessage
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
