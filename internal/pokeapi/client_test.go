package pokeapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alvaradoruben/myapp/internal/pokeapi"
)

// fixture returns a minimal valid PokeAPI /pokemon response for Pikachu.
func fixture() map[string]any {
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
			map[string]any{
				"base_stat": 35,
				"effort":    0,
				"stat":      map[string]any{"name": "hp", "url": ""},
			},
		},
		"types": []any{
			map[string]any{
				"slot": 1,
				"type": map[string]any{"name": "electric", "url": ""},
			},
		},
		"sprites": map[string]any{
			"front_default": "https://example.com/pikachu.png",
			"front_shiny":   "",
			"back_default":  "",
		},
	}
}

func TestGetPokemon_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pokemon/pikachu" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fixture())
	}))
	defer srv.Close()

	client := pokeapi.New(pokeapi.WithBaseURL(srv.URL))

	p, err := client.GetPokemon(context.Background(), "pikachu")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.ID != 25 {
		t.Errorf("expected ID 25, got %d", p.ID)
	}
	if p.Name != "pikachu" {
		t.Errorf("expected name 'pikachu', got %q", p.Name)
	}
	if len(p.Types) == 0 || p.Types[0] != "electric" {
		t.Errorf("expected type 'electric', got %v", p.Types)
	}
	if len(p.Abilities) == 0 || p.Abilities[0].Name != "static" {
		t.Errorf("expected ability 'static', got %v", p.Abilities)
	}
}

func TestGetPokemon_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	client := pokeapi.New(pokeapi.WithBaseURL(srv.URL))

	_, err := client.GetPokemon(context.Background(), "fakemon")
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if err != pokeapi.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetPokemon_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Deliberately slow — context should cancel first
		time.Sleep(500 * time.Millisecond)
		_ = json.NewEncoder(w).Encode(fixture())
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := pokeapi.New(pokeapi.WithBaseURL(srv.URL))

	_, err := client.GetPokemon(ctx, "pikachu")
	if err == nil {
		t.Fatal("expected context timeout error, got nil")
	}
}

func TestGetPokemon_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := pokeapi.New(pokeapi.WithBaseURL(srv.URL))

	_, err := client.GetPokemon(context.Background(), "pikachu")
	if err == nil {
		t.Fatal("expected error for 500 status, got nil")
	}
}

func TestListPokemon_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pokemon" {
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("expected limit=2, got %q", got)
		}
		if got := r.URL.Query().Get("offset"); got != "0" {
			t.Fatalf("expected offset=0, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"count":    1302,
			"next":     "https://pokeapi.co/api/v2/pokemon?offset=2&limit=2",
			"previous": nil,
			"results": []any{
				map[string]any{"name": "bulbasaur", "url": "https://pokeapi.co/api/v2/pokemon/1/"},
				map[string]any{"name": "ivysaur", "url": "https://pokeapi.co/api/v2/pokemon/2/"},
			},
		})
	}))
	defer srv.Close()

	client := pokeapi.New(pokeapi.WithBaseURL(srv.URL))

	list, err := client.ListPokemon(context.Background(), 2, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if list.Count != 1302 {
		t.Fatalf("expected count 1302, got %d", list.Count)
	}
	if len(list.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(list.Results))
	}
	if list.Results[0].ImageURL == "" {
		t.Fatal("expected derived image URL, got empty string")
	}
}
