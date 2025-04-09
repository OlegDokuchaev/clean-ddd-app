package middleware

import (
	"api-gateway/internal/infrastructure/logger"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

const maxBodyLogSize = 1024 // 1 KB

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// Write to the underlying ResponseWriter.
	n, err := w.ResponseWriter.Write(b)
	// Append the data to our buffer.
	w.body.Write(b)
	// If the buffer exceeds maxBodyLogSize, truncate it.
	if w.body.Len() > maxBodyLogSize {
		w.body.Truncate(maxBodyLogSize)
	}
	return n, err
}

func GinLoggingMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time to calculate latency.
		startTime := time.Now()

		// Replace the default ResponseWriter with our custom writer to capture the response body.
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = blw

		// Read the request body using a limited reader to avoid excessive memory usage.
		requestBody, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodyLogSize+1))
		if err != nil {
			log.Error("Error reading request body", map[string]any{"error": err})
		}
		// Truncate the request body if it exceeds the allowed maximum size.
		if len(requestBody) > maxBodyLogSize {
			requestBody = append(requestBody[:maxBodyLogSize], []byte("... (truncated)")...)
		}
		// Replace the request body so that it can be read again by the handlers.
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

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
			"request_body":     requestBody,
			"response_headers": c.Writer.Header(),
			"response_body":    blw.body.String(),
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
