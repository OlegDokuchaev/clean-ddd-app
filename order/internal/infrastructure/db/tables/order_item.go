package tables

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Count     int
	Price     decimal.Decimal `gorm:"type:numeric(10, 2)"`
}
