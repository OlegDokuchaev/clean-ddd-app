package commands

import (
	"context"
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

	CmdEnvelope struct {
		Ctx context.Context
		Cmd *CmdMessage
	}
)

type AssignCourierCmd struct {
	OrderID uuid.UUID
}
