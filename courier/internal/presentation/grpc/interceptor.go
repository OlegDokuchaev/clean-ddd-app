package courierv1

import (
	"context"
	"courier/internal/infrastructure/logger"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Log request start
		startTime := time.Now()
		method := path.Base(info.FullMethod)

		logger.Info("gRPC request received", map[string]any{
			"method":  method,
			"request": req,
		})

		// Call the request handler
		resp, err := handler(ctx, req)

		// Log request processing result
		duration := time.Since(startTime)
		statusCode := status.Code(err)

		logFields := map[string]any{
			"method":      method,
			"status_code": statusCode.String(),
			"duration_ms": duration.Milliseconds(),
		}

		if err != nil {
			logFields["error"] = err.Error()
			logger.Error("gRPC request failed with error", logFields)
		} else {
			logger.Info("gRPC request processed successfully", logFields)
		}

		return resp, err
	}
}

func StreamLoggingInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Log stream request start
		startTime := time.Now()
		method := path.Base(info.FullMethod)

		logger.Info("gRPC stream request started", map[string]any{
			"method":           method,
			"is_client_stream": info.IsClientStream,
			"is_server_stream": info.IsServerStream,
		})

		// Call the stream handler
		err := handler(srv, ss)

		// Log stream request processing result
		duration := time.Since(startTime)
		statusCode := status.Code(err)

		logFields := map[string]any{
			"method":      method,
			"status_code": statusCode.String(),
			"duration_ms": duration.Milliseconds(),
		}

		if err != nil {
			logFields["error"] = err.Error()
			logger.Error("gRPC stream request failed with error", logFields)
		} else {
			logger.Info("gRPC stream request processed successfully", logFields)
		}

		return err
	}
}
