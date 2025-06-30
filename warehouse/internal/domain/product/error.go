package product

import "errors"

var (
	ErrInvalidProductPrice = errors.New("invalid product price")
	ErrInvalidProductName  = errors.New("invalid product name")
)
