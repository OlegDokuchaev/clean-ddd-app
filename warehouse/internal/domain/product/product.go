package product

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ID      uuid.UUID
	Name    string
	Price   decimal.Decimal
	Created time.Time
	Image   Image
}
