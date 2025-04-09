package response

import (
	domainErrors "api-gateway/internal/domain/errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcToHTTP = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusRequestTimeout,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
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
