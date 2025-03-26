package commands

import (
	"encoding/json"
	"github.com/google/uuid"
)

const (
	AssignCourierCmdName CmdMessageName = "create_order.assign_courier"
)

type (
	CmdMessageName string

	CmdMessage struct {
		ID      uuid.UUID
		Name    CmdMessageName
		Payload json.RawMessage
	}
)

type AssignCourierCmd struct {
	OrderID uuid.UUID
}
