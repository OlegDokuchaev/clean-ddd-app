package outbox

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidOutboxMessage = errors.New("invalid outbox message")
)

func parseError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("outbox message not published: %w", err)
}
