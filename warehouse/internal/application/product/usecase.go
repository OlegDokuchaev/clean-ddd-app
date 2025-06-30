package product

import (
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, data CreateDto) (uuid.UUID, error)
}
