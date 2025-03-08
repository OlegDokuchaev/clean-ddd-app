package tables

import "github.com/google/uuid"

type OutboxMessage struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Type    string
	Payload []byte
}
