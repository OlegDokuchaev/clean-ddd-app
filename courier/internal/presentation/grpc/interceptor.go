package courierv1

import (
	"context"
	"courier/internal/infrastructure/logger"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func grpcLog(
	log logger.Logger,
	component string,
	action string,
	method string,
	level logger.Level,
	message string,
	extra map[string]any,
) {
	fields := map[string]any{
		"component": component,
		"action":    action,
		"method":    method,
	}
	for k, v := range extra {
		fields[k] = v
	}

	log.Log(level, message, fields)
}

func LoggingInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Log request start
		method := path.Base(info.FullMethod)

		grpcLog(log, "grpc_unary", "request_received", method, logger.Info,
			"gRPC request received", map[string]any{"request": req})

		startTime := time.Now()

		// Call the request handler
		resp, err := handler(ctx, req)

		// Log request processing result
		duration := time.Since(startTime)
		statusCode := status.Code(err)

		extra := map[string]any{
			"status_code": statusCode.String(),
			"duration_ms": duration.Milliseconds(),
			"response":    resp,
		}
		if err != nil {
			extra["error"] = err.Error()
			grpcLog(log, "grpc_unary", "request_failed", method, logger.Error,
				"gRPC request failed", extra)
		} else {
			grpcLog(log, "grpc_unary", "request_success", method, logger.Info,
				"gRPC request processed successfully", extra)
		}

		return resp, err
	}
}

func StreamLoggingInterceptor(log logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Log request start
		method := path.Base(info.FullMethod)

		extraStart := map[string]any{
			"is_client_stream": info.IsClientStream,
			"is_server_stream": info.IsServerStream,
		}
		grpcLog(log, "grpc_stream", "stream_started", method, logger.Info,
			"gRPC stream request started", extraStart)

		startTime := time.Now()

		// Call the request handler
		err := handler(srv, ss)

		// Log request processing result
		duration := time.Since(startTime)
		statusCode := status.Code(err)

		extraFinish := map[string]any{
			"status_code": statusCode.String(),
			"duration_ms": duration.Milliseconds(),
		}
		if err != nil {
			extraFinish["error"] = err.Error()
			grpcLog(log, "grpc_stream", "stream_failed", method, logger.Error,
				"gRPC stream request failed", extraFinish)
		} else {
			grpcLog(log, "grpc_stream", "stream_success", method, logger.Info,
				"gRPC stream request processed successfully", extraFinish)
		}

		return err
	}
}
