package product

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strings"
	"time"
	domain "warehouse/internal/domain/common"
)

func Create(name string, price decimal.Decimal) (*Product, []*domain.Event, error) {
	if price.LessThanOrEqual(decimal.Zero) {
		return nil, []*domain.Event{}, ErrInvalidProductPrice
	}
	if strings.TrimSpace(name) == "" {
		return nil, []*domain.Event{}, ErrInvalidProductName
	}

	product := &Product{
		ID:      uuid.New(),
		Name:    name,
		Price:   price,
		Created: time.Now(),
	}
	event := NewCreatedEvent(CreatedPayload{
		ProductID: product.ID,
	})

	return product, []*domain.Event{&event}, nil
}
