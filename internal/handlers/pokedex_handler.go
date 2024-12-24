package handlers

import (
	"encoding/json"

	"github.com/nguyensngoc108/pokemon-game/internal/repositories"
	"github.com/nguyensngoc108/pokemon-game/utils"
)

type PokeDexHandler struct {
	PokedexReposiotry *repositories.PokedexRepository
}

func NewPokeDexHandler(pokedexReposiotry *repositories.PokedexRepository) *PokeDexHandler {
	return &PokeDexHandler{
		PokedexReposiotry: pokedexReposiotry,
	}
}

func (s *PokeDexHandler) GetPokemon(name string) ([]byte, error) {
	data, err := s.PokedexReposiotry.GetMonsterByID(utils.PokeMap[name])
	if err != nil {
		return []byte{}, err
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return []byte{}, err
	}
	return jsonData, nil
}
