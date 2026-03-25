package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/pkg/logger"
)

func newTestHandler() *Handler {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "testapp",
			Env:     "test",
			Version: "0.0.0",
		},
	}
	log := logger.New("error", "test")
	return New(cfg, log, nil) // pokeClient is nil; unused in health handler tests
}

func TestHealthCheck(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var body healthResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", body.Status)
	}
}

func TestReadinessCheck(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	h.ReadinessCheck(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAppInfo(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/info", nil)
	w := httptest.NewRecorder()

	h.AppInfo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body appInfoResponse
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Name != "testapp" {
		t.Errorf("expected name 'testapp', got %q", body.Name)
	}
}
