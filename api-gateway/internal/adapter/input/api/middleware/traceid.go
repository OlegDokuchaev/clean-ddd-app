package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func TraceIDHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if span := trace.SpanFromContext(c.Request.Context()); span != nil {
			if sc := span.SpanContext(); sc.IsValid() {
				c.Header("trace-id", sc.TraceID().String())
			}
		}
	}
}
