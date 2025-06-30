package middleware

import (
	"api-gateway/internal/infrastructure/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func GinLoggingMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time to calculate latency.
		startTime := time.Now()

		// Process the request.
		c.Next()

		// Calculate latency.
		latency := time.Since(startTime)

		// Log detailed information about the HTTP request/response.
		fields := map[string]interface{}{
			"component":        "gin_http",
			"action":           "request_processed",
			"status":           c.Writer.Status(),
			"method":           c.Request.Method,
			"path":             c.Request.URL.Path,
			"query":            c.Request.URL.RawQuery,
			"client_ip":        c.ClientIP(),
			"user_agent":       c.Request.UserAgent(),
			"request_headers":  c.Request.Header,
			"request_size":     c.Request.ContentLength,
			"response_headers": c.Writer.Header(),
			"duration_ms":      latency.Milliseconds(),
		}

		if len(c.Errors.ByType(gin.ErrorTypePrivate)) > 0 {
			fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			log.Error("HTTP request failed", fields)
		} else {
			log.Info("HTTP request processed successfully", fields)
		}
	}
}
