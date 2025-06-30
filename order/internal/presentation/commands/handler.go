package commands

import (
	"context"
	"encoding/json"
	"fmt"
	createOrder "order/internal/application/order/saga/create_order"
	orderUsecase "order/internal/application/order/usecase"
	createOrderConsumer "order/internal/presentation/saga/create_order"
)

type Handler interface {
	Handle(ctx context.Context, cmdMsg *CmdMessage) (*createOrderConsumer.ResMessage, error)
}

type HandlerImpl struct {
	usecase orderUsecase.UseCase
}

func NewHandler(usecase orderUsecase.UseCase) *HandlerImpl {
	return &HandlerImpl{usecase: usecase}
}

func (h *HandlerImpl) Handle(ctx context.Context, cmdMsg *CmdMessage) (*createOrderConsumer.ResMessage, error) {
	switch cmdMsg.Name {
	case CancelOutOfStockCmdName:
		var cmd createOrder.CancelOutOfStockCmd
		if err := json.Unmarshal(cmdMsg.Payload, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse CancelOutOfStockCmd: %w", err)
		}
		return h.onCancelOutOfStock(ctx, cmd), nil

	case CancelCourierNotFoundCmdName:
		var cmd createOrder.CancelCourierNotFoundCmd
		if err := json.Unmarshal(cmdMsg.Payload, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse CancelCourierNotFoundCmd: %w", err)
		}
		return h.onCancelCourierNotFoundCmd(ctx, cmd), nil

	case BeginDeliveryCmdName:
		var cmd createOrder.BeginDeliveryCmd
		if err := json.Unmarshal(cmdMsg.Payload, &cmd); err != nil {
			return nil, fmt.Errorf("failed to parse BeginDeliveryCmd: %w", err)
		}
		return h.onBeginDelivery(ctx, cmd), nil
	}

	return nil, fmt.Errorf("unknown command: %s", cmdMsg.Name)
}

func (h *HandlerImpl) onCancelOutOfStock(
	ctx context.Context,
	cmd createOrder.CancelOutOfStockCmd,
) *createOrderConsumer.ResMessage {
	_ = h.usecase.CancelOutOfStock(ctx, cmd.OrderID)
	return nil
}

func (h *HandlerImpl) onCancelCourierNotFoundCmd(
	ctx context.Context,
	cmd createOrder.CancelCourierNotFoundCmd,
) *createOrderConsumer.ResMessage {
	_ = h.usecase.CancelCourierNotFound(ctx, cmd.OrderID)
	return nil
}

func (h *HandlerImpl) onBeginDelivery(
	ctx context.Context,
	cmd createOrder.BeginDeliveryCmd,
) *createOrderConsumer.ResMessage {
	data := orderUsecase.BeginDeliveryDto{
		OrderID:   cmd.OrderID,
		CourierID: cmd.CourierID,
	}
	_ = h.usecase.BeginDelivery(ctx, data)
	return nil
}

var _ Handler = (*HandlerImpl)(nil)
