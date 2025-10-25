package order

import (
	"time"

	"github.com/google/uuid"
)

type Delivery struct {
	CourierID *uuid.UUID
	Address   string
	Arrived   *time.Time
}
