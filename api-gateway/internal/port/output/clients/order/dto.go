package order

import (
	orderDto "api-gateway/internal/domain/dtos/order"
)

type CreateDto struct {
	CustomerID string
	Address    string
	Items      []orderDto.ItemDto
}
