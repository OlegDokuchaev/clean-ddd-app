package create_order

import "context"

type Publisher interface {
	PublishReserveItemsCmd(ctx context.Context, cmd ReserveItemsCmd) error
	PublishReleaseItemsCmd(ctx context.Context, cmd ReleaseItemsCmd) error
	PublishCancelOutOfStockCmd(ctx context.Context, cmd CancelOutOfStockCmd) error
	PublishAssignCourierCmd(ctx context.Context, cmd AssignCourierCmd) error
	PublishBeginDeliveryCmd(ctx context.Context, cmd BeginDeliveryCmd) error
	PublishCancelCourierNotFoundCmd(ctx context.Context, cmd CancelCourierNotFoundCmd) error
}
