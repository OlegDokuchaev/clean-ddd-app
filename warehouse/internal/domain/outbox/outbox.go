package outbox

import "github.com/google/uuid"

type Message struct {
	ID      uuid.UUID
	Name    string
	Payload []byte
}
