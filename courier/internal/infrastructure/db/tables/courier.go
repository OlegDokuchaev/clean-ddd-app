package tables

import (
	"github.com/google/uuid"
	"time"
)

type Courier struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string
	Phone    string `gorm:"unique"`
	Password []byte
	Created  time.Time
}
