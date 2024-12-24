package models

// Move represents a Pokémon move.
type Descriptions struct {
	Description string `json:"description"`
	ID          string `json:"_id"`
	Rev         string `json:"_rev"`
}
