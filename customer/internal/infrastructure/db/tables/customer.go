package tables

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Phone    string    `gorm:"unique"`
	Email    string
	Password []byte
	Created  time.Time
}
