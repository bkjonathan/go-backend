package app

import (
	"thomas-backend/internal/config"
	authDomain "thomas-backend/internal/domain/auth"
	userDomain "thomas-backend/internal/domain/user"
	"thomas-backend/internal/middleware"
	"thomas-backend/pkg/jwtutil"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App is the central dependency-injection container.
// All repositories, services, and handlers are constructed here.
// Adding a new domain means adding one field and a few lines in New().
type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB

	// Handlers — used by the router to register routes.
	AuthHandler *authDomain.Handler
	UserHandler *userDomain.Handler

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware
}

// New constructs the entire application dependency graph.
func New(cfg *config.Config, logger *zap.Logger, db *gorm.DB) *App {
	validate := validator.New()
	tokenManager := jwtutil.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	// --- Auth domain ---
	authRepo := authDomain.NewRepository(db)
	authService := authDomain.NewService(authRepo, validate, tokenManager, logger)
	authHandler := authDomain.NewHandler(authService, logger)

	// --- User domain ---
	userRepo := userDomain.NewRepository(db)
	userService := userDomain.NewService(userRepo, validate, logger)
	userHandler := userDomain.NewHandler(userService, logger)

	// --- Middleware ---
	authMiddleware := middleware.NewAuthMiddleware(tokenManager, logger)

	return &App{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
}
