package order

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Status     Status
	Created    time.Time
	Version    uuid.UUID
	Delivery   Delivery
	Items      []Item
}

func (o *Order) NoteCanceledByCustomer() error {
	switch o.Status {
	case Delivering:
		o.Status = CustomerCanceled
		return nil

	default:
		return ErrUnsupportedStatusTransition
	}
}

func (o *Order) NoteCanceledOutOfStock() error {
	switch o.Status {
	case Created:
		o.Status = CanceledOutOfStock
		return nil

	default:
		return ErrUnsupportedStatusTransition
	}
}

func (o *Order) NoteCanceledCourierNotFound() error {
	switch o.Status {
	case Created:
		o.Status = CanceledCourierNotFound
		return nil

	default:
		return ErrUnsupportedStatusTransition
	}
}

func (o *Order) NoteDelivering(CourierID uuid.UUID) error {
	switch o.Status {
	case Created:
		o.Status = Delivering
		o.Delivery.CourierID = &CourierID
		return nil

	default:
		return ErrUnsupportedStatusTransition
	}
}

func (o *Order) NoteDelivered() error {
	switch o.Status {
	case Delivering:
		now := time.Now()
		o.Status = Delivered
		o.Delivery.Arrived = &now
		return nil

	default:
		return ErrUnsupportedStatusTransition
	}
}
