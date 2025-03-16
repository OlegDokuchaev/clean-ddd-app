package commands

import "github.com/google/uuid"

type (
	ResMessageName    string
	ResMessagePayload interface{}

	ResMessage struct {
		ID      uuid.UUID
		Name    ResMessageName
		Payload ResMessagePayload
	}
)
