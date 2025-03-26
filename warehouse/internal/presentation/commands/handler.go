package commands

import (
	"context"
	"encoding/json"
	"fmt"
	itemApplication "warehouse/internal/application/item"
)

type Handler interface {
	Handle(ctx context.Context, cmdMsg *CmdMessage) (*ResMessage, error)
}

type HandlerImpl struct {
	usecase itemApplication.UseCase
}

func NewHandler(usecase itemApplication.UseCase) *HandlerImpl {
	return &HandlerImpl{usecase: usecase}
}

func (h *HandlerImpl) Handle(ctx context.Context, cmdMsg *CmdMessage) (*ResMessage, error) {
	switch cmdMsg.Name {
	case ReserveItemsCmdName:
		var cmd ReserveItemsCmd
		if err := json.Unmarshal(cmdMsg.Payload, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse ReserveItemsCmd: %w", err)
		}
		return h.onReserveItems(ctx, cmd), nil

	case ReleaseItemsCmdName:
		var cmd ReleaseItemsCmd
		if err := json.Unmarshal(cmdMsg.Payload, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse ReleaseItemsCmd: %w", err)
		}
		return h.onReleaseItems(ctx, cmd), nil
	}

	return nil, fmt.Errorf("unknown command: %s", cmdMsg.Name)
}

func (h *HandlerImpl) onReserveItems(ctx context.Context, cmd ReserveItemsCmd) *ResMessage {
	data := toReserveItemsDto(cmd)

	err := h.usecase.Reserve(ctx, data)

	if err != nil {
		return toItemsReservationFailed(cmd.OrderID)
	}
	return toItemsReserved(cmd.OrderID)
}

func (h *HandlerImpl) onReleaseItems(ctx context.Context, cmd ReleaseItemsCmd) *ResMessage {
	data := toReleaseItemsDto(cmd)

	err := h.usecase.Release(ctx, data)

	if err != nil {
		return nil
	}
	return toItemsReleased(cmd.OrderID)
}

var _ Handler = (*HandlerImpl)(nil)
