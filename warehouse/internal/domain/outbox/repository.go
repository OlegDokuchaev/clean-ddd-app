package outbox

import "context"

type Repository interface {
	Create(ctx context.Context, message *Message) error
	GetAll(ctx context.Context) ([]*Message, error)
	Delete(ctx context.Context, message *Message) error
}
