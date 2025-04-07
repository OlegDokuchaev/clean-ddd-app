package response

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/url"
)

func AddLocationHeader(c *gin.Context, path string) {
	if c.Writer.Header().Get("Location") != "" {
		return
	}

	fullURL, err := url.JoinPath(c.Request.URL.Path, path)
	if err != nil {
		return
	}

	c.Header("Location", fullURL)
}

func AddLocationHeaderWithID(c *gin.Context, id uuid.UUID) {
	AddLocationHeader(c, id.String())
}
