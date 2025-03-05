package order

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	ProductID uuid.UUID
	Price     decimal.Decimal
	Count     int
}
