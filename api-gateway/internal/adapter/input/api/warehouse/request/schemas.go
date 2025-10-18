package warehouse_request

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ItemInfoSchema struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Count     int       `json:"count" binding:"required"`
}

type GetAllItemsRequest struct {
	Limit  int `form:"limit" binding:"required"`
	Offset int `form:"offset" binding:"required"`
}

type CreateProductRequest struct {
	Name  string          `json:"name" binding:"required"`
	Price decimal.Decimal `json:"price" binding:"required"`
}

type ReserveItemsRequest struct {
	Items []*ItemInfoSchema `json:"items" binding:"required,min=1,dive"`
}

type ReleaseItemsRequest struct {
	Items []*ItemInfoSchema `json:"items" binding:"required,min=1,dive"`
}
