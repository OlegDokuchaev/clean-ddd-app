package events

import (
	"github.com/google/uuid"
)

type (
	EventName    string
	EventPayload interface{}

	Event struct {
		ID      uuid.UUID
		Name    EventName
		Payload EventPayload
	}
)

const (
	ProductCreatedName EventName = "product.ProductCreated"
)

type ProductCreatedEvent struct {
	ProductID uuid.UUID
}
