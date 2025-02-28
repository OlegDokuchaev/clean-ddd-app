package order

import (
	"github.com/google/uuid"
	"time"
)

func Create(CustomerID uuid.UUID, Address string, Items []Item) *Order {
	return &Order{
		ID:         uuid.New(),
		CustomerID: CustomerID,
		Status:     Created,
		Created:    time.Now(),
		Version:    uuid.New(),
		Delivery: Delivery{
			CourierID: nil,
			Address:   Address,
			Arrived:   nil,
		},
		Items: Items,
	}
}
