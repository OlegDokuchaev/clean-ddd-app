package order

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type CreateDto struct {
	Address string
	Items   []ItemDto
}

type OrderDto struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Status     Status
	Created    time.Time
	Version    uuid.UUID
	Delivery   DeliveryDto
	Items      []ItemDto
}

type ItemDto struct {
	ProductID uuid.UUID
	Price     decimal.Decimal
	Count     int
}

type DeliveryDto struct {
	CourierID *uuid.UUID
	Address   string
	Arrived   *time.Time
}
