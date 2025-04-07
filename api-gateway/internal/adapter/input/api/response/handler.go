package response

import (
	domainErrors "api-gateway/internal/domain/errors"
	"errors"

	"github.com/gin-gonic/gin"
)

type ErrorResponseDetail struct {
	Detail string `json:"detail"`
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var (
		appErr  *domainErrors.AppError
		respErr *ErrorResponse
	)

	switch {
	case errors.As(err, &appErr):
		respErr = &ErrorResponse{appErr.HTTPCode, appErr.Message}

	case errors.As(err, &respErr):

	default:
		respErr = ErrInternal
	}

	AbortWithError(c, ErrInternal)
}

func AbortWithError(c *gin.Context, err *ErrorResponse) {
	c.AbortWithStatusJSON(err.Status, ErrorResponseDetail{Detail: err.Detail})
}
