package models

// Move represents a Pokémon move.
type Move struct {
	TypeName    string      `json:"type_name"`
	Identifier  string      `json:"identifier"`
	Power       interface{} `json:"power"`
	PP          interface{} `json:"pp"`
	Accuracy    interface{} `json:"accuracy"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	ID          string      `json:"_id"`
	Rev         string      `json:"_rev"`
}
