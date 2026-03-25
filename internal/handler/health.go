package handler

import (
	"net/http"
	"time"
)

type healthResponse struct {
	Status  string `json:"status"`
	Time    string `json:"time"`
	Version string `json:"version"`
}

type appInfoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

// HealthCheck handles GET /health — liveness probe.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_ = writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Time:    time.Now().UTC().Format(time.RFC3339),
		Version: h.cfg.App.Version,
	})
}

// ReadinessCheck handles GET /ready — readiness probe.
// Add dependency checks (DB connections, caches, etc.) here before going live.
func (h *Handler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	_ = writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// AppInfo handles GET /api/v1/info — returns application metadata.
func (h *Handler) AppInfo(w http.ResponseWriter, r *http.Request) {
	_ = writeJSON(w, http.StatusOK, appInfoResponse{
		Name:    h.cfg.App.Name,
		Version: h.cfg.App.Version,
		Env:     h.cfg.App.Env,
	})
}
