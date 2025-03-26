package outbox

import (
	"encoding/json"
	"github.com/google/uuid"
)

type KafkaMessageValue struct {
	ID      uuid.UUID
	Name    string
	Payload json.RawMessage
}
