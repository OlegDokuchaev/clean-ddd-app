package common

import "github.com/google/uuid"

type (
	EventPayload interface{}

	Event interface {
		ID() uuid.UUID
		Name() string
		Payload() EventPayload
	}

	EventBase[T EventPayload] struct {
		id      uuid.UUID
		payload T
	}
)

func (e EventBase[T]) ID() uuid.UUID {
	return e.id
}
func (e EventBase[T]) Payload() EventPayload {
	return e.payload
}

func NewEvent[P EventPayload, T ~struct{ EventBase[P] }](payload P) T {
	return T{
		EventBase: EventBase[P]{
			id:      uuid.New(),
			payload: payload,
		},
	}
}
