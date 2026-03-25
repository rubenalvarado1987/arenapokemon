package handler

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/internal/server"
	"github.com/alvaradoruben/myapp/pkg/logger"
)

var (
	once       sync.Once
	appHandler http.Handler
	initErr    error
)

// Handler is the Vercel serverless entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		_ = godotenv.Load()

		log := logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("APP_ENV"))
		slog.SetDefault(log)

		cfg, err := config.Load()
		if err != nil {
			initErr = err
			return
		}

		appHandler = server.NewHandler(cfg, log)
	})

	if initErr != nil {
		http.Error(w, "failed to initialize application", http.StatusInternalServerError)
		return
	}

	appHandler.ServeHTTP(w, r)
}