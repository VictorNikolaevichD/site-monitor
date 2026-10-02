package interceptor

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const requestIDMetadataKey = "x-request-id"

type requestIDContextKey struct{}

func UnaryLogging(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		requestID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(requestIDMetadataKey); len(vals) > 0 {
				requestID = vals[0]
			}
		}
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx = context.WithValue(ctx, requestIDContextKey{}, requestID)

		resp, err := handler(ctx, req)

		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}

		logger.Info("grpc request",
			"request_id", requestID,
			"method", info.FullMethod,
			"code", code.String(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return resp, err
	}
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDContextKey{}).(string)
	return id
}
