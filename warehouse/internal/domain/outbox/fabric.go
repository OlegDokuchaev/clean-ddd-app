package outbox

import (
	"encoding/json"
	domain "warehouse/internal/domain/common"
)

func Create(event domain.Event) (*Message, error) {
	payload, err := parsePayload(event)
	if err != nil {
		return nil, err
	}
	return &Message{
		ID:      event.ID(),
		Type:    event.Name(),
		Payload: payload,
	}, nil
}

func parsePayload(event domain.Event) ([]byte, error) {
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return nil, err
	}
	return payload, nil
}
