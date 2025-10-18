package customer

import "errors"

var (
	ErrOtpInvalid          = errors.New("otp invalid")
	ErrOtpExpired          = errors.New("otp expired")
	ErrOtpAttemptsExceeded = errors.New("otp attempts exceeded")
)
