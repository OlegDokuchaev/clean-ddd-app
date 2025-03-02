package outbox

import "errors"

var (
	ErrInvalidOutboxPayload = errors.New("invalid outbox payload")
)
