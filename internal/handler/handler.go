package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/internal/pokeapi"
)

// Handler holds application-level dependencies shared across HTTP handlers.
type Handler struct {
	cfg        *config.Config
	log        *slog.Logger
	pokeClient *pokeapi.Client
}

// New creates a new Handler with the given dependencies.
func New(cfg *config.Config, log *slog.Logger, pokeClient *pokeapi.Client) *Handler {
	return &Handler{cfg: cfg, log: log, pokeClient: pokeClient}
}

// RegisterRoutes registers all application routes on the given mux.
// Uses Go 1.22+ enhanced ServeMux with method and path parameter support.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Frontend)

	// Health & readiness probes
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /ready", h.ReadinessCheck)

	// API v1
	mux.HandleFunc("GET /api/v1/info", h.AppInfo)
	mux.HandleFunc("GET /api/v1/pokemon", h.ListPokemon)
	mux.HandleFunc("GET /api/v1/pokemon/{name}", h.GetPokemon)
}

// writeJSON encodes data as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// writeError writes a structured JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	_ = writeJSON(w, status, map[string]string{"error": message})
}
