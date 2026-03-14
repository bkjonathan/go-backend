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
	httprouter "thomas-backend/internal/http/router"
	"time"

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

	//validate := validator.New()

	router := httprouter.New(&cfg, logger)

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
