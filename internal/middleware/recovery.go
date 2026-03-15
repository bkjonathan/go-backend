package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"thomas-backend/internal/apperror"
	"thomas-backend/internal/response"

	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error(
						"panic recovered",
						zap.Any("panic", recovered),
						zap.ByteString("stack", debug.Stack()),
						zap.String("path", r.URL.Path),
						zap.String("request_id", GetRequestID(r.Context())),
					)

					response.WriteFromError(w, apperror.Internal(
						"internal server error",
						fmt.Errorf("panic: %v", recovered),
					))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
