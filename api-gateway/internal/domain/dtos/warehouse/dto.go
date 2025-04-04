package warehouse

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ItemInfoDto struct {
	ProductID uuid.UUID
	Count     int
}

type ItemDto struct {
	ItemID  uuid.UUID
	Count   int
	Product ProductDto
	Version uuid.UUID
}

type ProductDto struct {
	ProductID uuid.UUID
	Name      string
	Price     decimal.Decimal
	Created   time.Time
}

type CreateProductDto struct {
	Name  string
	Price decimal.Decimal
}
