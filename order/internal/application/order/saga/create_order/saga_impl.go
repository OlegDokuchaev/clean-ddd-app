package create_order

import (
	"context"
	orderDomain "order/internal/domain/order"
)

type SagaImpl struct {
	publisher  Publisher
	repository orderDomain.Repository
}

func New(publisher Publisher, repository orderDomain.Repository) Saga {
	return &SagaImpl{
		publisher:  publisher,
		repository: repository,
	}
}

func (s *SagaImpl) HandleItemsReserved(ctx context.Context, event ItemsReserved) error {
	cmd := AssignCourierCmd(event)
	return s.publisher.PublishAssignCourierCmd(ctx, cmd)
}

func (s *SagaImpl) HandleItemsReservationFailed(ctx context.Context, event ItemsReservationFailed) error {
	cmd := CancelOutOfStockCmd(event)
	return s.publisher.PublishCancelOutOfStockCmd(ctx, cmd)
}

func (s *SagaImpl) HandleCourierAssignmentFailed(ctx context.Context, event CourierAssignmentFailed) error {
	order, err := s.repository.GetByID(ctx, event.OrderID)
	if err != nil {
		return err
	}

	orderItems := domainItemsToOrderItems(order.Items)

	cmd := ReleaseItemsCmd{
		OrderID: event.OrderID,
		Items:   orderItems,
	}
	return s.publisher.PublishReleaseItemsCmd(ctx, cmd)
}

func (s *SagaImpl) HandleItemsReleased(ctx context.Context, event ItemsReleased) error {
	cmd := CancelCourierNotFoundCmd(event)
	return s.publisher.PublishCancelCourierNotFoundCmd(ctx, cmd)
}

func (s *SagaImpl) HandleCourierAssigned(ctx context.Context, event CourierAssigned) error {
	cmd := BeginDeliveryCmd(event)
	return s.publisher.PublishBeginDeliveryCmd(ctx, cmd)
}

var _ Saga = (*SagaImpl)(nil)
