package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"thomas-backend/internal/config"
	"thomas-backend/internal/database"
	httprouter "thomas-backend/internal/http/router"
	"thomas-backend/internal/middleware"
	"thomas-backend/pkg/jwtutil"
	"time"

	authDomain "thomas-backend/internal/domain/auth"
	userDomain "thomas-backend/internal/domain/user"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("creating logger: %w", err))
	}
	defer func() {
		_ = logger.Sync()
	}()

	db, err := database.NewPostgresDB(ctx, cfg)
	if err != nil {
		logger.Fatal("connecting database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("getting sql db from gorm", zap.Error(err))
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	validate := validator.New()
	tokenManager := jwtutil.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	authRepo := authDomain.NewRepository(db)
	authService := authDomain.NewService(authRepo, validate, tokenManager, logger)
	authHandler := authDomain.NewHandler(authService, logger)

	userRepo := userDomain.NewRepository(db)
	userService := userDomain.NewService(userRepo, validate, logger)
	userHandler := userDomain.NewHandler(userService, logger)

	authMiddleware := middleware.NewAuthMiddleware(tokenManager, logger)

	router := httprouter.New(cfg, logger, authHandler, userHandler, authMiddleware)

	srv := &http.Server{
		Addr:              cfg.Server.Address,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("server started", zap.String("address", cfg.Server.Address), zap.String("env", cfg.App.Environment))
		if serveErr := srv.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			logger.Fatal("starting server", zap.Error(fmt.Errorf("listen and serve: %w", serveErr)))
		}
	}()
	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}

	logger.Info("shutting down")
}
