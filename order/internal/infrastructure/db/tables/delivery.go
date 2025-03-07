package tables

import (
	"time"

	"github.com/google/uuid"
)

type Delivery struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	OrderID   uuid.UUID `gorm:"uniqueIndex"`
	CourierID *uuid.UUID
	Address   string
	Arrived   *time.Time
}
