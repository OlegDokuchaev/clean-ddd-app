package response

import (
	domainErrors "api-gateway/internal/domain/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var grpcToHTTP = map[codes.Code]int{
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.NotFound:           http.StatusNotFound,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusConflict,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
}

var (
	ErrInternalServerError = domainErrors.NewAppError(http.StatusInternalServerError, "internal server error", nil)
)

func getHTTPCodeByStatus(st *status.Status) int {
	if code, ok := grpcToHTTP[st.Code()]; ok {
		return code
	}
	return http.StatusInternalServerError
}

func ParseGRPCError(err error) *domainErrors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return &domainErrors.AppError{
			HTTPCode: http.StatusInternalServerError,
			Message:  "internal server error",
			Err:      err,
		}
	}

	return &domainErrors.AppError{
		HTTPCode: getHTTPCodeByStatus(st),
		Message:  st.Message(),
		Err:      err,
	}
}
