package tables

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Phone    string `gorm:"unique"`
	Password []byte
	Created  time.Time

	FailedCount        int
	LockedUntil        *time.Time
	PasswordUpdated    time.Time
	MustChangePassword bool
}
