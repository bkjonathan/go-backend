package auth

import (
	"fmt"
	"net/http"
	"thomas-backend/internal/apperror"
	"thomas-backend/internal/common/httputil"
	"thomas-backend/internal/response"

	"go.uber.org/zap"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.WriteFromError(w, apperror.InvalidInput(
			"invalid request payload",
			fmt.Errorf("decoding register request: %w", err),
		))
		return
	}

	result, err := h.service.Register(r.Context(), req)
	if err != nil {
		h.logger.Warn("register request failed", zap.Error(fmt.Errorf("register service call: %w", err)))
		response.WriteFromError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, result, "user registered")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		response.WriteFromError(w, apperror.InvalidInput(
			"invalid request payload",
			fmt.Errorf("decoding login request: %w", err),
		))
		return
	}
	result, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Warn("login request failed", zap.Error(fmt.Errorf("login service call: %w", err)))
		response.WriteFromError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, result, "login successful")
}
