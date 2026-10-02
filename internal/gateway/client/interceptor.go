package client

import (
	"context"
	"log/slog"
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const requestIDMetadataKey = "x-request-id"

func UnaryClientLogging(logger *slog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		if id := middleware.RequestIDFromContext(ctx); id != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, requestIDMetadataKey, id)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)

		code := "OK"
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code().String()
			} else {
				code = "Unknown"
			}
		}

		logger.Info("grpc client call",
			"request_id", middleware.RequestIDFromContext(ctx),
			"method", method,
			"code", code,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}
}
