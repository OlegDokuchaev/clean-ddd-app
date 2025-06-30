package product

import (
	"strings"
	"time"
	domain "warehouse/internal/domain/common"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func Create(name string, price decimal.Decimal, imagePath string) (*Product, []domain.Event, error) {
	if price.LessThanOrEqual(decimal.Zero) {
		return nil, []domain.Event{}, ErrInvalidProductPrice
	}
	if strings.TrimSpace(name) == "" {
		return nil, []domain.Event{}, ErrInvalidProductName
	}

	product := &Product{
		ID:      uuid.New(),
		Name:    name,
		Price:   price,
		Created: time.Now(),
		Image: Image{
			Path: imagePath,
		},
	}
	event := domain.NewEvent[CreatedPayload, CreateEvent](CreatedPayload{
		ProductID: product.ID,
	})

	return product, []domain.Event{event}, nil
}
