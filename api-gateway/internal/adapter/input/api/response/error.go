package response

import (
	"net/http"
)

type ErrorResponse struct {
	Status int
	Detail string
}

func (e ErrorResponse) Error() string {
	return e.Detail
}

var (
	ErrInternal = &ErrorResponse{
		Status: http.StatusInternalServerError,
		Detail: "internal server error",
	}
	ErrInvalidId = &ErrorResponse{
		Status: http.StatusUnprocessableEntity,
		Detail: "invalid id",
	}
	ErrMissingAuthorization = &ErrorResponse{
		Status: http.StatusUnauthorized,
		Detail: "missing authorization",
	}
	ErrInvalidAuthorization = &ErrorResponse{
		Status: http.StatusUnauthorized,
		Detail: "invalid authorization",
	}
)
