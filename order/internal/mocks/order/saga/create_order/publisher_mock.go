package create_order

import (
	"context"
	createOrderSaga "order/internal/application/order/saga/create_order"

	"github.com/stretchr/testify/mock"
)

type PublisherMock struct {
	mock.Mock
}

func (p *PublisherMock) PublishReserveItemsCmd(ctx context.Context, cmd createOrderSaga.ReserveItemsCmd) error {
	args := p.Called(ctx, cmd)
	return args.Error(0)
}

func (p *PublisherMock) PublishReleaseItemsCmd(ctx context.Context, cmd createOrderSaga.ReleaseItemsCmd) error {
	args := p.Called(ctx, cmd)
	return args.Error(0)
}

func (p *PublisherMock) PublishCancelOutOfStockCmd(ctx context.Context, cmd createOrderSaga.CancelOutOfStockCmd) error {
	args := p.Called(ctx, cmd)
	return args.Error(0)
}

func (p *PublisherMock) PublishAssignCourierCmd(ctx context.Context, cmd createOrderSaga.AssignCourierCmd) error {
	args := p.Called(ctx, cmd)
	return args.Error(0)
}

func (p *PublisherMock) PublishBeginDeliveryCmd(ctx context.Context, cmd createOrderSaga.BeginDeliveryCmd) error {
	args := p.Called(ctx, cmd)
	return args.Error(0)
}

func (p *PublisherMock) PublishCancelCourierNotFoundCmd(ctx context.Context, cmd createOrderSaga.CancelCourierNotFoundCmd) error {
	args := p.Called(ctx, cmd)
	return args.Error(0)
}

var _ createOrderSaga.Publisher = (*PublisherMock)(nil)
