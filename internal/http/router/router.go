package router

import (
	"net/http"
	"thomas-backend/internal/app"
	"thomas-backend/internal/middleware"
	"thomas-backend/internal/response"

	"github.com/go-chi/chi/v5"
)

// New creates the HTTP router with all routes registered.
// It accepts the App container, making it trivial to add new domains.
func New(application *app.App) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery(application.Logger))
	r.Use(middleware.RequestLogger(application.Logger))
	r.Use(middleware.CORS(application.Config.CORS))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteSuccess(w, http.StatusOK, map[string]string{"status": "ok"}, "service healthy")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", application.AuthHandler.Register)
			r.Post("/login", application.AuthHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(application.AuthMiddleware.Handler)
			r.Route("/users", func(r chi.Router) {
				r.Get("/", application.UserHandler.List)
				r.Get("/me", application.UserHandler.GetMe)
				r.Get("/{id}", application.UserHandler.GetByID)
				r.Put("/{id}", application.UserHandler.Update)
				r.Delete("/{id}", application.UserHandler.Delete)
			})
		})
	})
	return r
}
