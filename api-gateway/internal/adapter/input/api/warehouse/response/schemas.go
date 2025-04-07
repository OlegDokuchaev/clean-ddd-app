package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ItemResponse struct {
	ItemID  uuid.UUID     `json:"item_id"`
	Count   int           `json:"count"`
	Product ProductSchema `json:"product"`
	Version string        `json:"version"`
}

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
}

type ProductSchema struct {
	ProductID uuid.UUID       `json:"product_id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Created   time.Time       `json:"created"`
}

type CreateProductResponse struct {
	ProductID uuid.UUID `json:"product_id"`
}
