package errors

type AppError struct {
	HTTPCode int
	Message  string
	Err      error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(httpCode int, message string, err error) *AppError {
	return &AppError{
		HTTPCode: httpCode,
		Message:  message,
		Err:      err,
	}
}
