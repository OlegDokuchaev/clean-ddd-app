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
	payloadBytes, err := json.Marshal(cmdMsg.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	switch cmdMsg.Name {
	case ItemsReservedName:
		var res createOrder.ItemsReserved
		if err = json.Unmarshal(payloadBytes, &res); err != nil {
			return fmt.Errorf("failed to parse ItemsReserved: %w", err)
		}
		return h.onItemsReserved(ctx, res)

	case ItemsReservationFailedName:
		var res createOrder.ItemsReservationFailed
		if err = json.Unmarshal(payloadBytes, &res); err != nil {
			return fmt.Errorf("failed to parse ItemsReservationFailed: %w", err)
		}
		return h.onItemsReservationFailed(ctx, res)

	case ItemsReleasedName:
		var res createOrder.ItemsReleased
		if err = json.Unmarshal(payloadBytes, &res); err != nil {
			return fmt.Errorf("failed to parse ItemsReleased: %w", err)
		}
		return h.onItemsReleased(ctx, res)

	case CourierAssignedName:
		var res createOrder.CourierAssigned
		if err = json.Unmarshal(payloadBytes, &res); err != nil {
			return fmt.Errorf("failed to parse CourierAssigned: %w", err)
		}
		return h.onCourierAssigned(ctx, res)

	case CourierAssignmentFailedName:
		var res createOrder.CourierAssignmentFailed
		if err = json.Unmarshal(payloadBytes, &res); err != nil {
			return fmt.Errorf("failed to parse CourierAssignmentFailed: %w", err)
		}
		return h.onCourierAssignmentFailed(ctx, res)
	}

	return fmt.Errorf("unknown command: %s", cmdMsg.Name)
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
