package order

import (
	orderDto "api-gateway/internal/domain/dtos/order"

	"github.com/google/uuid"
)

type CreateDto struct {
	CustomerID uuid.UUID
	Address    string
	Items      []orderDto.ItemDto
}
