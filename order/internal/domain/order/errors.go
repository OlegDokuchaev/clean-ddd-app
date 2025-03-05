package order

import "errors"

var (
	ErrUnsupportedStatusTransition = errors.New("unsupported order status transition")
	ErrInvalidAddress              = errors.New("invalid order address")
	ErrInvalidItems                = errors.New("invalid order items")
)
