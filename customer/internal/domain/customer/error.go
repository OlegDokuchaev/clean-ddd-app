package customer

import "errors"

var (
	ErrInvalidCustomerPassword = errors.New("invalid customer password")
	ErrInvalidCustomerPhone    = errors.New("invalid customer phone")
	ErrInvalidCustomerEmail    = errors.New("invalid customer email")
	ErrInvalidCustomerName     = errors.New("invalid customer name")
	ErrLocked                  = errors.New("customer is locked")
)
