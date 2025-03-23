package commands

import (
	"github.com/google/uuid"
)

const (
	AssignCourierCmdName CmdMessageName = "create_order.assign_courier"
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

type AssignCourierCmd struct {
	OrderID uuid.UUID
}
