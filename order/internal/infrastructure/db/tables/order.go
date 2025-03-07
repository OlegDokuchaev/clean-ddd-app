package tables

import (
	orderDomain "order/internal/domain/order"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	CustomerID uuid.UUID
	Status     orderDomain.Status `gorm:"type:order_status"`
	Created    time.Time
	Version    uuid.UUID
	Items      []OrderItem `gorm:"foreignKey:OrderID;"`
	Delivery   Delivery    `gorm:"foreignKey:OrderID;"`
}
