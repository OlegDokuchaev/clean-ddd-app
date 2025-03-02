package common

import "github.com/google/uuid"

type (
	EventPayload interface{}

	Event interface {
		ID() uuid.UUID
		Name() string
		Payload() EventPayload
	}

	event struct {
		id      uuid.UUID
		name    string
		payload EventPayload
	}
)

func (e event) ID() uuid.UUID {
	return e.id
}
func (e event) Name() string {
	return e.name
}
func (e event) Payload() EventPayload {
	return e.payload
}

func newEvent(name string, payload EventPayload) event {
	return event{
		id:      uuid.New(),
		name:    name,
		payload: payload,
	}
}

func NewEvent(name string, payload EventPayload) Event {
	return newEvent(name, payload)
}
