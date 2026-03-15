package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"thomas-backend/internal/apperror"
	"thomas-backend/internal/response"
	"thomas-backend/pkg/jwtutil"

	"go.uber.org/zap"
)

type AuthMiddleware struct {
	tokenManager *jwtutil.Manager
	logger       *zap.Logger
}

func NewAuthMiddleware(tokenManager *jwtutil.Manager, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := apperror.New(
				apperror.CodeUnauthorized,
				"missing authorization header",
				http.StatusUnauthorized,
				fmt.Errorf("authorization header is empty"),
			)
			response.WriteFromError(w, err)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			err := apperror.New(
				apperror.CodeUnauthorized,
				"invalid authorization header",
				http.StatusUnauthorized,
				fmt.Errorf("authorization format should be Bearer <token>"),
			)
			response.WriteFromError(w, err)
			return
		}

		claims, err := m.tokenManager.Verify(parts[1])
		if err != nil {
			m.logger.Warn("token verification failed", zap.Error(fmt.Errorf("verifying token: %w", err)))
			response.WriteFromError(w, apperror.New(
				apperror.CodeUnauthorized,
				"invalid token",
				http.StatusUnauthorized,
				fmt.Errorf("invalid jwt token: %w", err),
			))
			return
		}

		ctx := SetAuthUser(r.Context(), AuthenticatedUser{
			ID:    claims.UserID,
			Email: claims.Email,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
