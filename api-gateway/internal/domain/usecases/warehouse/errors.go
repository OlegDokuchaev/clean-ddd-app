package warehouse

import (
	domainErrors "api-gateway/internal/domain/errors"
	"net/http"
)

var (
	ErrUnauthorized = domainErrors.NewAppError(http.StatusUnauthorized, "unauthorized access", nil)
)
