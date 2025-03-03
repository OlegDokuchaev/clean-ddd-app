package item

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	itemDomain "warehouse/internal/domain/item"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, item *itemDomain.Item) error {
	args := r.Called(ctx, item)
	return args.Error(0)
}

func (r *RepositoryMock) Update(ctx context.Context, item *itemDomain.Item) error {
	args := r.Called(ctx, item)
	return args.Error(0)
}

func (r *RepositoryMock) GetByID(ctx context.Context, itemID uuid.UUID) (*itemDomain.Item, error) {
	args := r.Called(ctx, itemID)
	return args.Get(0).(*itemDomain.Item), args.Error(1)
}

func (r *RepositoryMock) GetAllByIDs(ctx context.Context, itemIDs ...uuid.UUID) ([]*itemDomain.Item, error) {
	argsForCalled := make([]interface{}, 0, len(itemIDs)+1)
	argsForCalled = append(argsForCalled, ctx)
	for _, id := range itemIDs {
		argsForCalled = append(argsForCalled, id)
	}

	args := r.Called(argsForCalled...)
	return args.Get(0).([]*itemDomain.Item), args.Error(1)
}

var _ itemDomain.Repository = (*RepositoryMock)(nil)
