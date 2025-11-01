package events

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
)

type (
	EventName string

	EventMessage struct {
		ID      uuid.UUID
		Name    EventName
		Payload json.RawMessage
	}

	EventEnvelope struct {
		Ctx       context.Context
		Msg       *EventMessage
		Topic     string
		Partition int
	}
)

const (
	ProductCreatedName EventName = "product.ProductCreated"
)

type ProductCreatedEvent struct {
	ProductID uuid.UUID
}
