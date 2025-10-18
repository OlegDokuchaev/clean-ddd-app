package otp_store

import "errors"

var (
	ErrOtpStoreUnavailable = errors.New("otp store unavailable")
	ErrOtpStoreCorrupted   = errors.New("otp store corrupted")
)
