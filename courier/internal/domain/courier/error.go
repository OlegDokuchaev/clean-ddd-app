package courier

import "errors"

var (
	ErrInvalidCourierPassword = errors.New("invalid courier password")
	ErrInvalidCourierPhone    = errors.New("invalid courier phone")
	ErrInvalidCourierName     = errors.New("invalid courier name")
)
