package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alvaradoruben/myapp/internal/config"
)

// Server wraps the standard http.Server.
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

// New creates and configures a new Server.
func New(cfg *config.Config, log *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Server.Addr(),
			Handler:      NewHandler(cfg, log),
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
