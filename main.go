package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nguyensngoc108/pokemon-game/internal/handlers"
	"github.com/nguyensngoc108/pokemon-game/internal/repositories"
	"github.com/nguyensngoc108/pokemon-game/utils"
)

func main() {

	pokedexRepository := repositories.NewPokedexRepository("../internal/models")
	pokedexHandler := handlers.NewPokeDexHandler(pokedexRepository)

	for key, value := range utils.PokeMap {
		pokemon, err := pokedexHandler.GetPokemon(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		filename := fmt.Sprintf("../internal/models/monsters/data/%s.json", value)
		//moveJSON, err := json.MarshalIndent(pokemon, "", "  ")
		//if err != nil {
		//	log.Printf("Failed to marshal move to JSON: %s\nError: %s", key, err)
		//	continue
		//}

		err = ioutil.WriteFile(filename, pokemon, 0644)
		if err != nil {
			log.Printf("Failed to write move to file: %s\nError: %s", filename, err)
			continue
		}

	}

}
