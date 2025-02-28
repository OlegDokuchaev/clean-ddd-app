package order

import "errors"

var (
	ErrUnsupportedStatusTransition = errors.New("unsupported order status transition")
)
