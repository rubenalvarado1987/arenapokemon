package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/internal/handler"
	"github.com/alvaradoruben/myapp/internal/middleware"
	"github.com/alvaradoruben/myapp/internal/pokeapi"
)

// Server wraps the standard http.Server.
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

// New creates and configures a new Server.
func New(cfg *config.Config, log *slog.Logger) *Server {
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

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Server.Addr(),
			Handler:      chain(mux),
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		log: log,
	}
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
