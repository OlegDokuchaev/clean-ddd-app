package commands

import (
	"context"
	"encoding/json"
	"fmt"
	createOrder "order/internal/application/order/saga/create_order"
	orderUsecase "order/internal/application/order/usecase"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"
	"order/internal/presentation/saga/create_order"
)

type Handler interface {
	Handle(ctx context.Context, cmdMsg *createOrderPublisher.CmdMessage) (*create_order.ResMessage, error)
}

type HandlerImpl struct {
	usecase orderUsecase.UseCase
}

func NewHandler(usecase orderUsecase.UseCase) *HandlerImpl {
	return &HandlerImpl{usecase: usecase}
}

func (h *HandlerImpl) Handle(ctx context.Context, cmdMsg *createOrderPublisher.CmdMessage) (*create_order.ResMessage, error) {
	payloadBytes, err := json.Marshal(cmdMsg.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	switch cmdMsg.Name {
	case createOrderPublisher.CancelOutOfStockCmdName:
		var cmd createOrder.CancelOutOfStockCmd
		if err = json.Unmarshal(payloadBytes, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse CancelOutOfStockCmd: %w", err)
		}
		return h.onCancelOutOfStock(ctx, cmd), nil

	case createOrderPublisher.CancelCourierNotFoundCmdName:
		var cmd createOrder.CancelCourierNotFoundCmd
		if err = json.Unmarshal(payloadBytes, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse CancelCourierNotFoundCmd: %w", err)
		}
		return h.onCancelCourierNotFoundCmd(ctx, cmd), nil

	case createOrderPublisher.BeginDeliveryCmdName:
		var cmd createOrder.BeginDeliveryCmd
		if err = json.Unmarshal(payloadBytes, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse BeginDeliveryCmd: %w", err)
		}
		return h.onBeginDelivery(ctx, cmd), nil
	}

	return nil, fmt.Errorf("unknown command: %s", cmdMsg.Name)
}

func (h *HandlerImpl) onCancelOutOfStock(ctx context.Context, cmd createOrder.CancelOutOfStockCmd) *create_order.ResMessage {
	_ = h.usecase.CancelOutOfStock(ctx, cmd.OrderID)
	return nil
}

func (h *HandlerImpl) onCancelCourierNotFoundCmd(ctx context.Context, cmd createOrder.CancelCourierNotFoundCmd) *create_order.ResMessage {
	_ = h.usecase.CancelCourierNotFound(ctx, cmd.OrderID)
	return nil
}

func (h *HandlerImpl) onBeginDelivery(ctx context.Context, cmd createOrder.BeginDeliveryCmd) *create_order.ResMessage {
	data := orderUsecase.BeginDeliveryDto{
		OrderID:   cmd.OrderID,
		CourierID: cmd.CourierID,
	}
	_ = h.usecase.BeginDelivery(ctx, data)
	return nil
}

var _ Handler = (*HandlerImpl)(nil)
