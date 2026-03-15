package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"thomas-backend/internal/apperror"

	"github.com/go-chi/chi/v5"
)

// DecodeJSON reads the request body as JSON into dst.
// Unknown fields are rejected.
func DecodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decoding json payload: %w", err)
	}
	return nil
}

// PathID extracts an int64 "id" path parameter from the request.
func PathID(r *http.Request) (int64, error) {
	return PathParamInt64(r, "id")
}

// PathParamInt64 extracts a named int64 path parameter from the request.
func PathParamInt64(r *http.Request, name string) (int64, error) {
	raw := chi.URLParam(r, name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apperror.InvalidInput(
			fmt.Sprintf("invalid %s parameter", name),
			fmt.Errorf("parsing %s parameter %q: %w", name, raw, err),
		)
	}
	return id, nil
}

// QueryIntOrDefault parses a query-string value as int, returning fallback on failure.
func QueryIntOrDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}
