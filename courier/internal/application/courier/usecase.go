package courier

import (
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	AssignOrder(ctx context.Context, orderID uuid.UUID) (uuid.UUID, error)
}
