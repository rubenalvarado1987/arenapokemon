package pokeapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "https://pokeapi.co/api/v2"

// ErrNotFound is returned when a Pokémon does not exist in the API.
var ErrNotFound = errors.New("pokemon not found")

// Client is an HTTP client for the PokéAPI.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithTimeout overrides the default HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithBaseURL overrides the API base URL (useful for testing).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// New returns a configured PokéAPI Client.
// A dedicated http.Client is used (never the package-level http.DefaultClient)
// to honour timeouts and avoid shared transport state.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// GetPokemon fetches a single Pokémon by name or numeric ID string.
// It returns ErrNotFound when the API responds with 404.
func (c *Client) GetPokemon(ctx context.Context, nameOrID string) (*Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, nameOrID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// handled below
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status %d from pokeapi", resp.StatusCode)
	}

	var raw pokemonResponse
	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return mapPokemon(&raw), nil
}

// ListPokemon fetches a paginated Pokemon listing from the PokéAPI.
func (c *Client) ListPokemon(ctx context.Context, limit, offset int) (*PokemonList, error) {
	url := fmt.Sprintf("%s/pokemon?limit=%d&offset=%d", c.baseURL, limit, offset)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from pokeapi", resp.StatusCode)
	}

	var raw pokemonListResponse
	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	items := make([]PokemonListItem, len(raw.Results))
	for i, result := range raw.Results {
		items[i] = PokemonListItem{
			Name:     result.Name,
			ImageURL: spriteURLFromResourceURL(result.URL),
		}
	}

	return &PokemonList{
		Count:   raw.Count,
		Limit:   limit,
		Offset:  offset,
		Results: items,
	}, nil
}

// mapPokemon converts the raw PokeAPI response into our public model.
func mapPokemon(r *pokemonResponse) *Pokemon {
	abilities := make([]Ability, len(r.Abilities))
	for i, a := range r.Abilities {
		abilities[i] = Ability{
			Name:     a.Ability.Name,
			IsHidden: a.IsHidden,
			Slot:     a.Slot,
		}
	}

	stats := make([]Stat, len(r.Stats))
	for i, s := range r.Stats {
		stats[i] = Stat{
			Name:     s.Stat.Name,
			BaseStat: s.BaseStat,
			Effort:   s.Effort,
		}
	}

	types := make([]string, len(r.Types))
	for i, t := range r.Types {
		types[i] = t.Type.Name
	}

	return &Pokemon{
		ID:             r.ID,
		Name:           r.Name,
		BaseExperience: r.BaseExperience,
		Height:         r.Height,
		Weight:         r.Weight,
		Abilities:      abilities,
		Stats:          stats,
		Types:          types,
		Sprites: Sprites{
			FrontDefault: r.Sprites.FrontDefault,
			FrontShiny:   r.Sprites.FrontShiny,
			BackDefault:  r.Sprites.BackDefault,
		},
	}
}

func spriteURLFromResourceURL(resourceURL string) string {
	trimmed := strings.TrimSuffix(resourceURL, "/")
	lastSlash := strings.LastIndex(trimmed, "/")
	if lastSlash == -1 || lastSlash == len(trimmed)-1 {
		return ""
	}
	id := trimmed[lastSlash+1:]
	if _, err := strconv.Atoi(id); err != nil {
		return ""
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%s.png", id)
}
