package request

import (
	"api-gateway/internal/adapter/input/api/response"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func ParseParamUUID(c *gin.Context, key string) (uuid.UUID, error) {
	if id, err := uuid.Parse(c.Param(key)); err != nil {
		return uuid.Nil, response.ErrInvalidId
	} else {
		return id, nil
	}
}

func ParseInput(c *gin.Context, input any, b binding.Binding) error {
	if err := c.ShouldBindWith(input, b); err != nil {
		return &response.ErrorResponse{Status: http.StatusUnprocessableEntity, Detail: err.Error()}
	}
	return nil
}

func ParseBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", response.ErrMissingAuthorization
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", response.ErrInvalidAuthorization
	}

	return parts[1], nil
}

func ParseAccessToken(c *gin.Context) (string, error) {
	accessToken := c.GetHeader("X-Access-Token")
	if accessToken != "" {
		return accessToken, nil
	}

	return "", response.ErrMissingAuthorization
}

func ParseFormFile(c *gin.Context, name string) (*multipart.FileHeader, error) {
	fileHeader, err := c.FormFile(name)
	if err != nil {
		return nil, response.ErrMissingFormFile
	}
	return fileHeader, nil
}
