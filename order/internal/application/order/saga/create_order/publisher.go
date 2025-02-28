package create_order

type Publisher interface {
	PublishReserveItemsCmd(cmd ReserveItemsCmd) error
	PublishReleaseItemsCmd(cmd ReleaseItemsCmd) error
	PublishCancelOutOfStockCmd(cmd CancelOutOfStockCmd) error
	PublishAssignCourierCmd(cmd AssignCourierCmd) error
	PublishBeginDeliveryCmd(cmd BeginDeliveryCmd) error
	PublishCancelCourierNotFoundCmd(cmd CancelCourierNotFoundCmd) error
}
