package tables

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OutboxMessage struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string
	Payload  []byte
	Metadata datatypes.JSON
}
