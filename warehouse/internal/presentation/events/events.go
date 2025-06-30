package events

import (
	"encoding/json"
	"github.com/google/uuid"
)

type (
	EventName string

	Event struct {
		ID      uuid.UUID
		Name    EventName
		Payload json.RawMessage
	}
)

const (
	ProductCreatedName EventName = "product.ProductCreated"
)

type ProductCreatedEvent struct {
	ProductID uuid.UUID
}
