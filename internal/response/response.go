package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"thomas-backend/internal/apperror"
)

type SuccessBody struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorBody struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func WriteSuccess(w http.ResponseWriter, status int, data interface{}, message string) {
	writeJSON(w, status, SuccessBody{
		Data:    data,
		Message: message,
	})
}

func WriteError(w http.ResponseWriter, status int, message, code string) {
	writeJSON(w, status, ErrorBody{
		Error: message,
		Code:  code,
	})
}

func WriteFromError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		WriteError(w, appErr.Status, appErr.Message, appErr.Code)
		return
	}

	WriteError(w, http.StatusInternalServerError, "internal server error", apperror.CodeInternal)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("encoding response: %v", err), http.StatusInternalServerError)
	}
}
