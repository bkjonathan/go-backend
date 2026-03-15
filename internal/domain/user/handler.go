package user

import (
	"fmt"
	"net/http"
	"thomas-backend/internal/apperror"
	"thomas-backend/internal/common/httputil"
	"thomas-backend/internal/middleware"
	"thomas-backend/internal/response"

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
	limit := int32(httputil.QueryIntOrDefault(r.URL.Query().Get("limit"), 100))
	offset := int32(httputil.QueryIntOrDefault(r.URL.Query().Get("offset"), 0))

	users, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("list users failed", zap.Error(fmt.Errorf("list users service call: %w", err)))
		response.WriteFromError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, users, "users fetched")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.PathID(r)
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
		response.WriteFromError(w, apperror.Unauthorized(
			"unauthenticated",
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
	id, err := httputil.PathID(r)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	var req UpdateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.WriteFromError(w, apperror.InvalidInput(
			"invalid request payload",
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
	id, err := httputil.PathID(r)
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
