package commands

import (
	"context"
	courierApplication "courier/internal/application/courier"
	"encoding/json"
	"fmt"
)

type Handler interface {
	Handle(ctx context.Context, cmdMsg *CmdMessage) (*ResMessage, error)
}

type HandlerImpl struct {
	usecase courierApplication.UseCase
}

func NewHandler(usecase courierApplication.UseCase) *HandlerImpl {
	return &HandlerImpl{usecase: usecase}
}

func (h *HandlerImpl) Handle(ctx context.Context, cmdMsg *CmdMessage) (*ResMessage, error) {
	payloadBytes, err := json.Marshal(cmdMsg.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	switch cmdMsg.Name {
	case AssignCourierCmdName:
		var cmd AssignCourierCmd
		if err = json.Unmarshal(payloadBytes, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse ReserveItemsCmd: %w", err)
		}
		return h.onAssignOrder(ctx, cmd), nil
	}

	return nil, fmt.Errorf("unknown command: %s", cmdMsg.Name)
}

func (h *HandlerImpl) onAssignOrder(ctx context.Context, cmd AssignCourierCmd) *ResMessage {
	courierID, err := h.usecase.AssignOrder(ctx, cmd.OrderID)

	if err != nil {
		return toCourierAssignmentFailed(cmd.OrderID)
	}
	return toCourierAssigned(cmd.OrderID, courierID)
}

var _ Handler = (*HandlerImpl)(nil)
