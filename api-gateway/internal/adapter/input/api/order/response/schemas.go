package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderResponse struct {
	ID         uuid.UUID      `json:"id"`
	CustomerID uuid.UUID      `json:"customer_id"`
	Status     string         `json:"status"`
	Created    time.Time      `json:"created"`
	Version    string         `json:"version"`
	Delivery   DeliverySchema `json:"delivery"`
	Items      []ItemSchema   `json:"items"`
}

type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}

type DeliverySchema struct {
	CourierID *uuid.UUID `json:"courier_id,omitempty"`
	Address   string     `json:"address"`
	Arrived   *time.Time `json:"arrived,omitempty"`
}

type ItemSchema struct {
	ProductID uuid.UUID       `json:"product_id"`
	Price     decimal.Decimal `json:"price"`
	Count     int             `json:"count"`
}
