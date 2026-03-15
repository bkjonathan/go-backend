package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"thomas-backend/internal/apperror"
	"thomas-backend/internal/middleware"
	"thomas-backend/internal/response"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit := int32(parseIntOrDefault(r.URL.Query().Get("limit"), 100))
	offset := int32(parseIntOrDefault(r.URL.Query().Get("offset"), 0))

	users, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("list users failed", zap.Error(fmt.Errorf("list users service call: %w", err)))
		response.WriteFromError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, users, "users fetched")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	user, getErr := h.service.GetByID(r.Context(), id)
	if getErr != nil {
		response.WriteFromError(w, getErr)
		return
	}

	response.WriteSuccess(w, http.StatusOK, user, "user fetched")
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.WriteFromError(w, apperror.New(
			apperror.CodeUnauthorized,
			"unauthenticated",
			http.StatusUnauthorized,
			fmt.Errorf("auth user missing in context"),
		))
		return
	}

	user, err := h.service.GetMe(r.Context(), authUser.ID)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, user, "current user fetched")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	var req UpdateRequest
	if err := decodeJSON(r, &req); err != nil {
		response.WriteFromError(w, apperror.New(
			apperror.CodeInvalidInput,
			"invalid request payload",
			http.StatusBadRequest,
			fmt.Errorf("decoding update payload: %w", err),
		))
		return
	}

	updated, updateErr := h.service.Update(r.Context(), id, req)
	if updateErr != nil {
		response.WriteFromError(w, updateErr)
		return
	}

	response.WriteSuccess(w, http.StatusOK, updated, "user updated")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.WriteFromError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, nil, "user deleted")
}

func pathID(r *http.Request) (int64, error) {
	idRaw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		return 0, apperror.New(
			apperror.CodeInvalidInput,
			"invalid id parameter",
			http.StatusBadRequest,
			fmt.Errorf("parsing id parameter %q: %w", idRaw, err),
		)
	}
	return id, nil
}

func parseIntOrDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func decodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decoding json payload: %w", err)
	}
	return nil
}
