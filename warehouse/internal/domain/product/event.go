package product

import (
	"github.com/google/uuid"
	domain "warehouse/internal/domain/common"
)

const (
	CreatedEventName = "product.ProductCreated"
)

type CreatedPayload struct {
	ProductID uuid.UUID
}

func NewCreatedEvent(payload CreatedPayload) domain.Event {
	return domain.NewEvent(CreatedEventName, payload)
}
