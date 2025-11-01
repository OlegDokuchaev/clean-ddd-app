package builders

import (
	orderDomain "order/internal/domain/order"
	"time"

	"github.com/google/uuid"
)

type OrderBuilder struct {
	id         uuid.UUID
	customerID uuid.UUID
	status     orderDomain.Status
	created    time.Time
	version    uuid.UUID
	delivery   orderDomain.Delivery
	items      []orderDomain.Item
}

func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{
		id:         uuid.New(),
		customerID: uuid.New(),
		status:     orderDomain.Created,
		created:    time.Now(),
		version:    uuid.New(),
		delivery: orderDomain.Delivery{
			Address: "Default address",
		},
		items: []orderDomain.Item{},
	}
}

func (b *OrderBuilder) WithCustomerID(id uuid.UUID) *OrderBuilder {
	b.customerID = id
	return b
}

func (b *OrderBuilder) WithStatus(status orderDomain.Status) *OrderBuilder {
	b.status = status
	return b
}

func (b *OrderBuilder) WithAddress(addr string) *OrderBuilder {
	b.delivery.Address = addr
	return b
}

func (b *OrderBuilder) WithDelivery(delivery orderDomain.Delivery) *OrderBuilder {
	b.delivery = delivery
	return b
}

func (b *OrderBuilder) WithItems(items []orderDomain.Item) *OrderBuilder {
	b.items = items
	return b
}

func (b *OrderBuilder) Build() *orderDomain.Order {
	return &orderDomain.Order{
		ID:         b.id,
		CustomerID: b.customerID,
		Status:     b.status,
		Created:    b.created,
		Version:    b.version,
		Delivery:   b.delivery,
		Items:      b.items,
	}
}
