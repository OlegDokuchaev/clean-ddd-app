package create_order

import (
	"context"
	"encoding/json"
	"fmt"
	createOrder "order/internal/application/order/saga/create_order"
)

type Handler interface {
	Handle(ctx context.Context, cmdMsg *ResMessage) error
}

type HandlerImpl struct {
	saga createOrder.Saga
}

func NewHandler(saga createOrder.Saga) *HandlerImpl {
	return &HandlerImpl{saga: saga}
}

func (h *HandlerImpl) Handle(ctx context.Context, cmdMsg *ResMessage) error {
	switch cmdMsg.Name {
	case ItemsReservedName:
		return h.handleItemsReserved(ctx, cmdMsg)
	case ItemsReservationFailedName:
		return h.handleItemsReservationFailed(ctx, cmdMsg)
	case ItemsReleasedName:
		return h.handleItemsReleased(ctx, cmdMsg)
	case CourierAssignedName:
		return h.handleCourierAssigned(ctx, cmdMsg)
	case CourierAssignmentFailedName:
		return h.handleCourierAssignmentFailed(ctx, cmdMsg)
	default:
		return fmt.Errorf("unknown command: %s", cmdMsg.Name)
	}
}

func (h *HandlerImpl) handleItemsReserved(ctx context.Context, cmdMsg *ResMessage) error {
	if cmdMsg.Name != ItemsReservedName {
		return fmt.Errorf("unexpected command: %s", cmdMsg.Name)
	}

	var res createOrder.ItemsReserved
	if err := json.Unmarshal(cmdMsg.Payload, &res); err != nil {
		return fmt.Errorf("failed to parse ItemsReserved: %w", err)
	}

	return h.onItemsReserved(ctx, res)
}

func (h *HandlerImpl) handleItemsReservationFailed(ctx context.Context, cmdMsg *ResMessage) error {
	if cmdMsg.Name != ItemsReservationFailedName {
		return fmt.Errorf("unexpected command: %s", cmdMsg.Name)
	}

	var res createOrder.ItemsReservationFailed
	if err := json.Unmarshal(cmdMsg.Payload, &res); err != nil {
		return fmt.Errorf("failed to parse ItemsReservationFailed: %w", err)
	}

	return h.onItemsReservationFailed(ctx, res)
}

func (h *HandlerImpl) handleItemsReleased(ctx context.Context, cmdMsg *ResMessage) error {
	if cmdMsg.Name != ItemsReleasedName {
		return fmt.Errorf("unexpected command: %s", cmdMsg.Name)
	}

	var res createOrder.ItemsReleased
	if err := json.Unmarshal(cmdMsg.Payload, &res); err != nil {
		return fmt.Errorf("failed to parse ItemsReleased: %w", err)
	}

	return h.onItemsReleased(ctx, res)
}

func (h *HandlerImpl) handleCourierAssigned(ctx context.Context, cmdMsg *ResMessage) error {
	if cmdMsg.Name != CourierAssignedName {
		return fmt.Errorf("unexpected command: %s", cmdMsg.Name)
	}

	var res createOrder.CourierAssigned
	if err := json.Unmarshal(cmdMsg.Payload, &res); err != nil {
		return fmt.Errorf("failed to parse CourierAssigned: %w", err)
	}

	return h.onCourierAssigned(ctx, res)
}

func (h *HandlerImpl) handleCourierAssignmentFailed(ctx context.Context, cmdMsg *ResMessage) error {
	if cmdMsg.Name != CourierAssignmentFailedName {
		return fmt.Errorf("unexpected command: %s", cmdMsg.Name)
	}

	var res createOrder.CourierAssignmentFailed
	if err := json.Unmarshal(cmdMsg.Payload, &res); err != nil {
		return fmt.Errorf("failed to parse CourierAssignmentFailed: %w", err)
	}

	return h.onCourierAssignmentFailed(ctx, res)
}

func (h *HandlerImpl) onItemsReserved(ctx context.Context, res createOrder.ItemsReserved) error {
	return h.saga.HandleItemsReserved(ctx, res)
}

func (h *HandlerImpl) onItemsReservationFailed(ctx context.Context, res createOrder.ItemsReservationFailed) error {
	return h.saga.HandleItemsReservationFailed(ctx, res)
}

func (h *HandlerImpl) onItemsReleased(ctx context.Context, res createOrder.ItemsReleased) error {
	return h.saga.HandleItemsReleased(ctx, res)
}

func (h *HandlerImpl) onCourierAssigned(ctx context.Context, res createOrder.CourierAssigned) error {
	return h.saga.HandleCourierAssigned(ctx, res)
}

func (h *HandlerImpl) onCourierAssignmentFailed(ctx context.Context, res createOrder.CourierAssignmentFailed) error {
	return h.saga.HandleCourierAssignmentFailed(ctx, res)
}

var _ Handler = (*HandlerImpl)(nil)
