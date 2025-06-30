package tables

import "github.com/google/uuid"

type OutboxMessage struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Name    string
	Payload []byte
}
