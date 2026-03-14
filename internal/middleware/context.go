package middleware

import "context"

type contextKey string

const (
	requestIDContextKey contextKey = "request_id"
	authUserContextKey  contextKey = "auth_user"
)

type AuthenticatedUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	value, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}
	return value
}

func SetAuthUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authUserContextKey, user)
}

func GetAuthUser(ctx context.Context) (AuthenticatedUser, bool) {
	value, ok := ctx.Value(authUserContextKey).(AuthenticatedUser)
	if !ok {
		return AuthenticatedUser{}, false
	}
	return value, true
}
