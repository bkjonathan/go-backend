package router

import (
	"net/http"
	"thomas-backend/internal/config"
	authDomain "thomas-backend/internal/domain/auth"
	userDomain "thomas-backend/internal/domain/user"
	"thomas-backend/internal/middleware"
	"thomas-backend/internal/response"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(
	cfg *config.Config,
	logger *zap.Logger,
	authHandler *authDomain.Handler,
	userHandler *userDomain.Handler,
	authMiddleware *middleware.AuthMiddleware,
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

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Handler)
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Get("/me", userHandler.GetMe)
				r.Get("/{id}", userHandler.GetByID)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})
		})
	})
	return r
}
