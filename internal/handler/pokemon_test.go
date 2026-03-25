package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvaradoruben/myapp/internal/config"
	"github.com/alvaradoruben/myapp/internal/handler"
	"github.com/alvaradoruben/myapp/internal/pokeapi"
	"github.com/alvaradoruben/myapp/pkg/logger"
)

// pokemonFixture is a minimal valid PokeAPI /pokemon response.
func pokemonFixture() map[string]any {
	return map[string]any{
		"id":              25,
		"name":            "pikachu",
		"base_experience": 112,
		"height":          4,
		"weight":          60,
		"abilities": []any{
			map[string]any{
				"ability":   map[string]any{"name": "static", "url": ""},
				"is_hidden": false,
				"slot":      1,
			},
		},
		"stats": []any{
			map[string]any{"base_stat": 35, "effort": 0, "stat": map[string]any{"name": "hp", "url": ""}},
		},
		"types": []any{
			map[string]any{"slot": 1, "type": map[string]any{"name": "electric", "url": ""}},
		},
		"sprites": map[string]any{
			"front_default": "https://example.com/pikachu.png",
			"front_shiny":   "",
			"back_default":  "",
		},
	}
}

func newHandlerWithMockPokeAPI(t *testing.T, upstream *httptest.Server) *handler.Handler {
	t.Helper()
	cfg := &config.Config{
		App: config.AppConfig{Name: "testapp", Env: "test", Version: "0.0.0"},
	}
	log := logger.New("error", "test")
	client := pokeapi.New(pokeapi.WithBaseURL(upstream.URL))
	return handler.New(cfg, log, client)
}

func TestGetPokemon_OK(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pokemonFixture())
	}))
	defer upstream.Close()

	h := newHandlerWithMockPokeAPI(t, upstream)

	// Use a mux so PathValue works correctly
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/pokemon/{name}", h.GetPokemon)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon/pikachu", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var p pokeapi.Pokemon
	if err := json.NewDecoder(w.Body).Decode(&p); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if p.Name != "pikachu" {
		t.Errorf("expected 'pikachu', got %q", p.Name)
	}
	if len(p.Types) == 0 || p.Types[0] != "electric" {
		t.Errorf("expected type 'electric', got %v", p.Types)
	}
}

func TestGetPokemon_NotFound(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer upstream.Close()

	h := newHandlerWithMockPokeAPI(t, upstream)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/pokemon/{name}", h.GetPokemon)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon/fakemon", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetPokemon_InvalidName(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should never be reached for truly invalid names
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	h := newHandlerWithMockPokeAPI(t, upstream)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/pokemon/{name}", h.GetPokemon)

	cases := []struct {
		name           string
		acceptedCodes  []int
	}{
		// Path traversal: Go's ServeMux cleans the path and issues a 307 redirect
		// before our handler runs — the attempt never reaches us.
		{"../etc/passwd", []int{http.StatusMovedPermanently, http.StatusPermanentRedirect,
			http.StatusFound, http.StatusTemporaryRedirect, http.StatusNotFound, http.StatusBadRequest}},
		// Our regex rejects uppercase and special characters with 400.
		{"Pikachu!", []int{http.StatusBadRequest}},
		{"PIKACHU", []int{http.StatusBadRequest}},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon/"+tc.name, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		ok := false
		for _, code := range tc.acceptedCodes {
			if w.Code == code {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("name %q: got unexpected status %d", tc.name, w.Code)
		}
	}
}

func TestListPokemon_OK(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pokemon" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count": 2,
			"next": nil,
			"previous": nil,
			"results": []any{
				map[string]any{"name": "bulbasaur", "url": "https://pokeapi.co/api/v2/pokemon/1/"},
				map[string]any{"name": "ivysaur", "url": "https://pokeapi.co/api/v2/pokemon/2/"},
			},
		})
	}))
	defer upstream.Close()

	h := newHandlerWithMockPokeAPI(t, upstream)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/pokemon", h.ListPokemon)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pokemon?limit=2&offset=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var list pokeapi.PokemonList
	if err := json.NewDecoder(w.Body).Decode(&list); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(list.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(list.Results))
	}
}

func TestFrontend_OK(t *testing.T) {
	h := newHandlerWithMockPokeAPI(t, httptest.NewServer(http.NotFoundHandler()))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Frontend)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if contentType := w.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected html content type, got %q", contentType)
	}
}
