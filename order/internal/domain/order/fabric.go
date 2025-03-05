package order

import (
	"time"

	"github.com/google/uuid"
)

func Create(CustomerID uuid.UUID, Address string, Items []Item) (*Order, error) {
	if !validateAddress(Address) {
		return nil, ErrInvalidAddress
	}
	if !validateItems(Items) {
		return nil, ErrInvalidItems
	}

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
	}, nil
}
