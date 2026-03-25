package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/internal/handler"
	"github.com/alvaradoruben/myapp/internal/middleware"
	"github.com/alvaradoruben/myapp/internal/pokeapi"
)

// NewHandler builds the application's HTTP handler chain.
func NewHandler(cfg *config.Config, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	pokeClient := pokeapi.New(
		pokeapi.WithTimeout(10 * time.Second),
	)

	h := handler.New(cfg, log, pokeClient)
	h.RegisterRoutes(mux)

	chain := middleware.Chain(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recover(log),
	)

	return chain(mux)
}