package customer

import "errors"

var (
	ErrInvalidCustomerPassword = errors.New("invalid customer password")
	ErrInvalidCustomerPhone    = errors.New("invalid customer phone")
	ErrInvalidCustomerName     = errors.New("invalid customer name")
)
