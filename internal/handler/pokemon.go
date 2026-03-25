package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/alvaradoruben/myapp/internal/pokeapi"
)

// pokemonNameRE matches valid Pokémon names/IDs: lowercase letters, digits, hyphens.
// This prevents path traversal or injection before the value is used in a URL.
var pokemonNameRE = regexp.MustCompile(`^[a-z0-9-]+$`)

// GetPokemon handles GET /api/v1/pokemon/{name}
//
// Path parameters:
//   - name: Pokémon name (e.g. "pikachu") or numeric ID (e.g. "25")
//
// Responses:
//   - 200 OK        — Pokémon data as JSON
//   - 400 Bad Request — invalid name format
//   - 404 Not Found  — Pokémon does not exist in PokéAPI
//   - 502 Bad Gateway — upstream PokéAPI error
func (h *Handler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if !pokemonNameRE.MatchString(name) {
		slog.WarnContext(r.Context(), "invalid pokemon name requested", "name", name)
		writeError(w, http.StatusBadRequest, "invalid pokemon name: use lowercase letters, digits and hyphens only")
		return
	}

	slog.InfoContext(r.Context(), "fetching pokemon", "name", name)

	pokemon, err := h.pokeClient.GetPokemon(r.Context(), name)
	if err != nil {
		if errors.Is(err, pokeapi.ErrNotFound) {
			writeError(w, http.StatusNotFound, "pokemon not found")
			return
		}
		slog.ErrorContext(r.Context(), "pokeapi request failed", "name", name, "error", err)
		writeError(w, http.StatusBadGateway, "could not retrieve pokemon data")
		return
	}

	_ = writeJSON(w, http.StatusOK, pokemon)
}

// ListPokemon handles GET /api/v1/pokemon — returns paginated Pokemon cards data.
func (h *Handler) ListPokemon(w http.ResponseWriter, r *http.Request) {
	limit := 24
	offset := 0

	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit < 1 || parsedLimit > 151 {
			writeError(w, http.StatusBadRequest, "invalid limit: use a number between 1 and 151")
			return
		}
		limit = parsedLimit
	}

	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil || parsedOffset < 0 {
			writeError(w, http.StatusBadRequest, "invalid offset: use a number greater than or equal to 0")
			return
		}
		offset = parsedOffset
	}

	list, err := h.pokeClient.ListPokemon(r.Context(), limit, offset)
	if err != nil {
		slog.ErrorContext(r.Context(), "pokeapi list request failed", "limit", limit, "offset", offset, "error", err)
		writeError(w, http.StatusBadGateway, "could not retrieve pokemon list")
		return
	}

	_ = writeJSON(w, http.StatusOK, list)
}
