package outbox

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	outboxDomain "warehouse/internal/domain/outbox"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, message *outboxDomain.Message) error {
	args := r.Called(ctx, message)
	return args.Error(0)
}

func (r *RepositoryMock) GetByID(ctx context.Context, messageID uuid.UUID) (*outboxDomain.Message, error) {
	args := r.Called(ctx, messageID)
	return args.Get(0).(*outboxDomain.Message), args.Error(1)
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]*outboxDomain.Message, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*outboxDomain.Message), args.Error(1)
}

func (r *RepositoryMock) Delete(ctx context.Context, message *outboxDomain.Message) error {
	args := r.Called(ctx, message)
	return args.Error(0)
}

var _ outboxDomain.Repository = (*RepositoryMock)(nil)
