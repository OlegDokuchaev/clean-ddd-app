package request

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	Address string        `json:"address" binding:"required"`
	Items   []*ItemSchema `json:"items" binding:"required,min=1,dive"`
}

type ItemSchema struct {
	ProductID uuid.UUID       `json:"product_id" binding:"required"`
	Price     decimal.Decimal `json:"price" binding:"required"`
	Count     int             `json:"count" binding:"required"`
}
