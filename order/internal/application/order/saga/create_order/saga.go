package create_order

import "context"

type Saga interface {
	HandleItemsReserved(ctx context.Context, event ItemsReserved) error
	HandleItemsReservationFailed(ctx context.Context, event ItemsReservationFailed) error
	HandleCourierAssignmentFailed(ctx context.Context, event CourierAssignmentFailed) error
	HandleCourierAssigned(ctx context.Context, event CourierAssigned) error
	HandleItemsReleased(ctx context.Context, event ItemsReleased) error
}
