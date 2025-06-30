package outbox

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, messageID uuid.UUID) (*Message, error)
	GetAll(ctx context.Context) ([]*Message, error)
	Delete(ctx context.Context, message *Message) error
}
