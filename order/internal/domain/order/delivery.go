package order

import (
	"github.com/google/uuid"
	"time"
)

type Delivery struct {
	CourierID *uuid.UUID
	Address   string
	Arrived   *time.Time
}
