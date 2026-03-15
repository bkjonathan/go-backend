package router

import (
	"net/http"
	"thomas-backend/internal/config"
	authDomain "thomas-backend/internal/domain/auth"
	"thomas-backend/internal/middleware"
	"thomas-backend/internal/response"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(
	cfg *config.Config,
	logger *zap.Logger,
	authHandler *authDomain.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS(cfg.CORS))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteSuccess(w, http.StatusOK, map[string]string{"status": "ok"}, "service healthy")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})
	})
	return r
}
