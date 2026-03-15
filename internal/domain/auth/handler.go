package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"thomas-backend/internal/apperror"
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
	if err := decodeJSON(r, &req); err != nil {
		response.WriteFromError(w, apperror.New(
			apperror.CodeInvalidInput,
			"invalid request payload",
			http.StatusBadRequest,
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
	if err := decodeJSON(r, &req); err != nil {
		response.WriteFromError(w, apperror.New(
			apperror.CodeInvalidInput,
			"invalid request payload",
			http.StatusBadRequest,
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
func decodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decoding json payload: %w", err)
	}
	return nil
}
