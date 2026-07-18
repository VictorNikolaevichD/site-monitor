package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
)

func Recovery(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				requestID := RequestIDFromContext(r.Context())

				logger.Error("panic", "request_id", requestID, "method", r.Method, "path", r.URL.Path, "error", rec, "stack", string(debug.Stack()))

				handler.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})

}
