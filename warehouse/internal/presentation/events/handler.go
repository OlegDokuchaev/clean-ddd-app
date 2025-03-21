package events

import (
	"context"
	"encoding/json"
	"fmt"
	"warehouse/internal/application/item"
)

type Handler interface {
	Handle(ctx context.Context, event *Event) error
}

type HandlerImpl struct {
	itemUseCase item.UseCase
}

func NewHandler(itemUseCase item.UseCase) *HandlerImpl {
	return &HandlerImpl{itemUseCase: itemUseCase}
}

func (h *HandlerImpl) Handle(ctx context.Context, event *Event) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	switch event.Name {
	case ProductCreatedName:
		var eventPayload ProductCreatedEvent
		if err = json.Unmarshal(payloadBytes, &eventPayload); err != nil {
			return fmt.Errorf("failed to parse ProductCreatedEvent: %w", err)
		}
		return h.onProductCreated(ctx, eventPayload)

	default:
		return fmt.Errorf("unknown event: %s", event.Name)
	}
}

func (h *HandlerImpl) onProductCreated(ctx context.Context, event ProductCreatedEvent) error {
	data := item.CreateDto{
		ProductID: event.ProductID,
	}

	_, err := h.itemUseCase.Create(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

var _ Handler = (*HandlerImpl)(nil)
