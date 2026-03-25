package pokeapi

// --- PokeAPI raw response types ---------------------------------------------

type pokemonResponse struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Abilities      []ability  `json:"abilities"`
	Stats          []stat     `json:"stats"`
	Types          []typeSlot `json:"types"`
	Sprites        sprites    `json:"sprites"`
}

type ability struct {
	Ability struct {
		Name string `json:"name"`
	} `json:"ability"`
	IsHidden bool `json:"is_hidden"`
	Slot     int  `json:"slot"`
}

type stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type typeSlot struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type sprites struct {
	FrontDefault     string `json:"front_default"`
	FrontShiny       string `json:"front_shiny"`
	BackDefault      string `json:"back_default"`
	FrontDefaultFemale string `json:"front_female"`
}

type pokemonListResponse struct {
	Count   int                 `json:"count"`
	Next    string              `json:"next"`
	Previous string             `json:"previous"`
	Results []pokemonListResult `json:"results"`
}

type pokemonListResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// --- Public API response types ----------------------------------------------

// Pokemon is the response model exposed by our API.
type Pokemon struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	BaseExperience int       `json:"base_experience"`
	Height         int       `json:"height"`       // in decimetres
	Weight         int       `json:"weight"`       // in hectograms
	Abilities      []Ability `json:"abilities"`
	Stats          []Stat    `json:"stats"`
	Types          []string  `json:"types"`
	Sprites        Sprites   `json:"sprites"`
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
	Count   int               `json:"count"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	Results []PokemonListItem `json:"results"`
}
