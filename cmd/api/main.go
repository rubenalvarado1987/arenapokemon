package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/internal/server"
	"github.com/alvaradoruben/myapp/pkg/logger"
)

func main() {
	// Load .env file if present (ignored in production)
	_ = godotenv.Load()

	// Initialize structured logger
	log := logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("APP_ENV"))
	slog.SetDefault(log)

	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create and configure HTTP server
	srv := server.New(cfg, log)

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server started", "addr", cfg.Server.Addr(), "env", cfg.App.Env)
		serverErr <- srv.Start()
	}()

	// Wait for OS signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-serverErr:
		slog.Error("server error", "error", err)
	case sig := <-quit:
		slog.Info("shutdown signal received", "signal", sig)
	}

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited gracefully")
}
