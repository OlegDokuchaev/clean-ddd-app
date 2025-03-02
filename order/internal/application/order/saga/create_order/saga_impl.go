package create_order

import (
	"context"
	"github.com/google/uuid"
)

type SagaImpl struct {
	publisher Publisher
}

func New(publisher Publisher) Saga {
	return &SagaImpl{
		publisher: publisher,
	}
}

func (s *SagaImpl) HandleItemsReserved(ctx context.Context, event ItemsReserved) error {
	cmd := AssignCourierCmd{
		OrderID: event.OrderID,
		Cmd: Cmd{
			ID: uuid.New(),
		},
	}
	return s.publisher.PublishAssignCourierCmd(ctx, cmd)
}

func (s *SagaImpl) HandleItemsReservationFailed(ctx context.Context, event ItemsReservationFailed) error {
	cmd := CancelOutOfStockCmd{
		OrderID: event.OrderID,
		Cmd: Cmd{
			ID: uuid.New(),
		},
	}
	return s.publisher.PublishCancelOutOfStockCmd(ctx, cmd)
}

func (s *SagaImpl) HandleCourierAssignmentFailed(ctx context.Context, event CourierAssignmentFailed) error {
	cmd := ReleaseItemsCmd{
		OrderID: event.OrderID,
		Cmd: Cmd{
			ID: uuid.New(),
		},
	}
	return s.publisher.PublishReleaseItemsCmd(ctx, cmd)
}

func (s *SagaImpl) HandleItemsReleased(ctx context.Context, event ItemsReleased) error {
	cmd := CancelCourierNotFoundCmd{
		OrderID: event.OrderID,
		Cmd: Cmd{
			ID: uuid.New(),
		},
	}
	return s.publisher.PublishCancelCourierNotFoundCmd(ctx, cmd)
}

func (s *SagaImpl) HandleCourierAssigned(ctx context.Context, event CourierAssigned) error {
	cmd := BeginDeliveryCmd{
		OrderID:   event.OrderID,
		CourierID: event.CourierID,
		Cmd: Cmd{
			ID: uuid.New(),
		},
	}
	return s.publisher.PublishBeginDeliveryCmd(ctx, cmd)
}

var _ Saga = (*SagaImpl)(nil)
