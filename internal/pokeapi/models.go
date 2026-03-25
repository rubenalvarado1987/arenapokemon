package pokeapi

// --- PokeAPI raw response types ---------------------------------------------

type pokemonResponse struct {
	Sprites        sprites    `json:"sprites"`
	Name           string     `json:"name"`
	Abilities      []ability  `json:"abilities"`
	Stats          []stat     `json:"stats"`
	Types          []typeSlot `json:"types"`
	ID             int        `json:"id"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
}

type ability struct {
	Ability struct {
		Name string `json:"name"`
	} `json:"ability"`
	IsHidden bool `json:"is_hidden"`
	Slot     int  `json:"slot"`
}

type stat struct {
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
}

type typeSlot struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
	Slot int `json:"slot"`
}

type sprites struct {
	FrontDefault       string `json:"front_default"`
	FrontShiny         string `json:"front_shiny"`
	BackDefault        string `json:"back_default"`
	FrontDefaultFemale string `json:"front_female"`
}

type pokemonListResponse struct {
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []pokemonListResult `json:"results"`
	Count    int                 `json:"count"`
}

type pokemonListResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// --- Public API response types ----------------------------------------------

// Pokemon is the response model exposed by our API.
type Pokemon struct {
	Sprites        Sprites   `json:"sprites"`
	Name           string    `json:"name"`
	Abilities      []Ability `json:"abilities"`
	Stats          []Stat    `json:"stats"`
	Types          []string  `json:"types"`
	ID             int       `json:"id"`
	BaseExperience int       `json:"base_experience"`
	Height         int       `json:"height"`
	Weight         int       `json:"weight"`
}

// Ability represents a Pokémon ability.
type Ability struct {
	Name     string `json:"name"`
	IsHidden bool   `json:"is_hidden"`
	Slot     int    `json:"slot"`
}

// Stat represents a base stat value.
type Stat struct {
	Name     string `json:"name"`
	BaseStat int    `json:"base_stat"`
	Effort   int    `json:"effort"`
}

// Sprites holds image URLs for the Pokémon.
type Sprites struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
	BackDefault  string `json:"back_default"`
}

// PokemonListItem represents a Pokemon entry in paginated listings.
type PokemonListItem struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// PokemonList is the paginated response model exposed by our API.
type PokemonList struct {
	Results []PokemonListItem `json:"results"`
	Count   int               `json:"count"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
}
