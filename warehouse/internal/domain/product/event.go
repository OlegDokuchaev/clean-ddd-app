package product

import (
	"github.com/google/uuid"
	domain "warehouse/internal/domain/common"
)

const (
	CreatedEventName = "product.ProductCreated"
)

type CreateEvent struct {
	domain.EventBase[CreatedPayload]
}

func (e CreateEvent) Name() string {
	return CreatedEventName
}

type CreatedPayload struct {
	ProductID uuid.UUID
}

var _ domain.Event = (*CreateEvent)(nil)
