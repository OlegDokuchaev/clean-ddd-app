package tables

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string
	Phone    string `gorm:"unique"`
	Password []byte
	Created  time.Time
}
