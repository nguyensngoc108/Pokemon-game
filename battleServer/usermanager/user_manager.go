package usermanager

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strings"

	"github.com/nguyensngoc108/pokemon-game/internal/models"
	"github.com/nguyensngoc108/pokemon-game/utils"
)

type User struct {
	Username      string
	Pokemons      []string
	Conn          net.Conn
	PokemonData   []*PokemonData
	ActivePokemon *PokemonData
	ActiveHP      int
}

type PokemonData struct {
	Monster             *models.Monster             `json:"monster"`
	Description         []*models.Descriptions      `json:"description"`
	Evolution           *models.Evolution           `json:"evolution"`
	Types               []*models.Types             `json:"types"`
	MonsterSupplemental *models.MonsterSupplemental `json:"monster_supplemental"`
	MonsterMoves        []*models.Move              `json:"monster_moves"`
}

type UserManager struct {
	Users         map[string]*User
	CurrentTurn   string
	BattleStarted bool
}

var userManagerInstance *UserManager

func GetUserManagerInstance() *UserManager {
	return userManagerInstance
}

func SetUserManagerInstance(um *UserManager) *UserManager {
	userManagerInstance = um
	return userManagerInstance
}

func (um *UserManager) GetUserPokemon(username string) string {
	return um.Users[username].ActivePokemon.Monster.Name
}
func (um *UserManager) GetUserPokemonHP(username string) int {
	return um.Users[username].ActivePokemon.Monster.HP
}
func (um *UserManager) GetUserPokemonActiveHP(username string) int {
	return um.Users[username].ActiveHP
}

func NewUserManager() *UserManager {
	return &UserManager{
		Users: make(map[string]*User),
	}
}

func (um *UserManager) AddUser(username string, conn net.Conn) *User {
	if _, exists := um.Users[username]; !exists {
		user := &User{
			Username: username,
			Pokemons: make([]string, 3),
			Conn:     conn,
		}
		um.Users[username] = user
		return user
	}
	return um.Users[username]
}

func (um *UserManager) UpdatePokemons(username, pokemon string, index int) {
	user, exists := um.Users[username]
	if exists {
		user.Pokemons[index-1] = pokemon
	}
}

func (um *UserManager) GetOpponentPokemons(username string) []string {
	for _, user := range um.Users {
		if user.Username != username {
			return user.Pokemons
		}
	}
	return nil
}

func (um *UserManager) AllPokemonsProvided() bool {
	// Check if there are exactly 2 connected players
	if len(um.Users) != 2 {
		return false
	}

	// Check if each player has provided 3 Pokemon
	for _, user := range um.Users {
		if len(user.Pokemons) != 3 {
			return false
		}
		for _, pokemon := range user.Pokemons {
			if pokemon == "" {
				return false
			}
		}
	}
	return true
}

func (um *UserManager) StartBattle() {
	// Broadcast the Pokemon information to both players
	um.broadcastPokemons()

	// Determine the player who goes first based on the speed of their first Pokemon
	um.determineTurnOrder()

	// Set initial Pokemon and HP
	for _, user := range um.Users {
		user.ActivePokemon = user.PokemonData[0]
		user.ActiveHP = user.ActivePokemon.Monster.HP
	}
	currentUser := um.Users[um.CurrentTurn]
	message := fmt.Sprintf("\n %s is at %d/%d HP. Choose your next move (type in 'normal' or 'special'): ", currentUser.ActivePokemon.Monster.Name, currentUser.ActiveHP, currentUser.ActivePokemon.Monster.HP)
	um.sendMessageToUser(currentUser, message)
}

func (um *UserManager) determineTurnOrder() {
	// Determine the player who goes first based on the speed of their first Pokemon
	user1 := um.Users[um.getPlayerNames()[0]]
	user2 := um.Users[um.getPlayerNames()[1]]

	if user1.PokemonData[0].Monster.Speed > user2.PokemonData[0].Monster.Speed {
		um.CurrentTurn = user1.Username
	} else {
		um.CurrentTurn = user2.Username
	}

	um.BattleStarted = true
}

func (um *UserManager) getPlayerNames() []string {
	var playerNames []string
	for username := range um.Users {
		playerNames = append(playerNames, username)
	}
	return playerNames
}

func (um *UserManager) broadcastPokemons() {
	for _, user := range um.Users {
		message := fmt.Sprintf("%s's Pokemon: %s\nOpponent's Pokemon: %s", user.Username, strings.Join(user.Pokemons, ", "), strings.Join(um.getOpponentPokemons(user), ", "))
		_, err := user.Conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", user.Username, err)
		}
		message = fmt.Sprintf("Battle Commence")
		_, err = user.Conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", user.Username, err)
		}
	}
}

func (um *UserManager) getOpponentPokemons(currentUser *User) []string {
	for _, user := range um.Users {
		if user != currentUser {
			return user.Pokemons
		}
	}
	return nil
}

func (um *UserManager) PerformBattle(moveType string, username string) {
	currentUser := um.Users[um.CurrentTurn]
	opponent := um.getOpponent(currentUser.Username)
	//message := fmt.Sprintf("\n %s is at %d/%d HP. Choose your next move (type in 'normal' or 'special'): ", currentUser.ActivePokemon.Monster.Name, currentUser.ActiveHP, currentUser.ActivePokemon.Monster.HP)
	//um.sendMessageToUser(currentUser, message)

	if username != currentUser.Username {
		return
	}

	// Verify that the move type is valid (either "normal attack" or "special attack")
	if moveType != "normal" && moveType != "special" && moveType != "quit" {
		um.sendMessageToUser(currentUser, fmt.Sprintf("Invalid move type: %s", moveType))
		return
	}
	if moveType == "quit" {
		um.announceWinner(opponent.Username)
		return
	}
	// Select a random move of the chosen type
	choosenMove := um.selectRandomMove(currentUser, moveType)

	// Calculate and apply the damage
	um.calculateAndApplyDamage(currentUser, opponent, choosenMove)

	// Check if the opponent's active Pokemon is knocked out
	if opponent.ActiveHP <= 0 {
		if um.hasLost(opponent) {
			um.announceWinner(currentUser.Username)
			return
		}
		um.replaceKnockedOutPokemon(opponent)
	}

	// Switch the turn to the other player
	um.switchTurn()
	message := fmt.Sprintf("\n %s is at %d/%d HP. Choose your next move (type in 'normal' or 'special' or 'quit'): ", opponent.ActivePokemon.Monster.Name, opponent.ActiveHP, opponent.ActivePokemon.Monster.HP)
	um.sendMessageToUser(opponent, message)
}

func (um *UserManager) selectRandomMove(user *User, moveType string) string {
	// Filter the moves that match the specified type
	var matchingMoves []*models.Move
	if moveType == "special" {
		for _, move := range user.ActivePokemon.MonsterMoves {
			if move.TypeName != "normal" && move.Power != "" {
				matchingMoves = append(matchingMoves, move)
			}
		}
	} else {
		for _, move := range user.ActivePokemon.MonsterMoves {
			if move.TypeName == moveType {
				matchingMoves = append(matchingMoves, move)
			}
		}
	}

	// If there are no matching moves, return an empty string
	if len(matchingMoves) == 0 {
		return ""
	}

	// Select a random move from the matching moves
	randomIndex := rand.Intn(len(matchingMoves))
	selectedMove := matchingMoves[randomIndex]
	return selectedMove.Name
}

func (um *UserManager) calculateAndApplyDamage(currentUser, defender *User, moveName string) {
	// Find the move object based on the move name
	var attackingMove *models.Move
	for _, move := range currentUser.ActivePokemon.MonsterMoves {
		if move.Name == moveName {
			attackingMove = move
			break
		}
	}

	if attackingMove == nil {
		um.sendMessageToUser(currentUser, fmt.Sprintf("%s does not have the move %s", currentUser.ActivePokemon.Monster.Name, moveName))
		return
	}

	// Calculate the damage based on the move type
	var damage int
	if attackingMove.TypeName == "normal" {
		damage = int(math.Abs(float64(currentUser.ActivePokemon.Monster.Attack - defender.ActivePokemon.Monster.Defense)))
	} else {
		damage = int(math.Abs(float64(currentUser.ActivePokemon.Monster.SpAtk - defender.ActivePokemon.Monster.SpDef)))
	}

	// Apply the damage to the defender's HP
	defender.ActiveHP = defender.ActiveHP - damage
	if defender.ActiveHP < 0 {
		defender.ActiveHP = 0
	}

	// Send the damage update to both players
	um.sendMessageToUser(currentUser, fmt.Sprintf("%s's %s did %d damage to %s's %s. %s's HP is now %d.", currentUser.Username, attackingMove.Name, damage, defender.Username, defender.ActivePokemon.Monster.Name, defender.Username, defender.ActiveHP))
	um.sendMessageToUser(defender, fmt.Sprintf("%s's %s did %d damage to %s's %s. %s's HP is now %d.", currentUser.Username, attackingMove.Name, damage, defender.Username, defender.ActivePokemon.Monster.Name, defender.Username, defender.ActiveHP))
}

func (um *UserManager) sendMessageToUser(user *User, message string) {
	_, err := user.Conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error sending message to %s: %v\n", user.Username, err)
	}
}

func (um *UserManager) getOpponent(username string) *User {
	for _, user := range um.Users {
		if user.Username != username {
			return user
		}
	}
	return nil
}

func (um *UserManager) replaceKnockedOutPokemon(user *User) {
	// Find the index of the knocked out Pokemon in the user's Pokemon list
	var knockedOutIndex int
	for i, pokemon := range user.Pokemons {
		if pokemon == user.ActivePokemon.Monster.Name {
			knockedOutIndex = i
			break
		}
	}

	// Select the next available Pokemon
	var nextPokemonIndex int
	for i, _ := range user.Pokemons {
		if i > knockedOutIndex && knockedOutIndex != 2 {
			nextPokemonIndex = i
			break
		}
	}

	// Update the user's active Pokemon and HP
	user.ActivePokemon = user.PokemonData[nextPokemonIndex]
	user.ActiveHP = user.ActivePokemon.Monster.HP
	// Check if the user has no more Pokemon left
	if knockedOutIndex == 2 {
		um.announceWinner(um.getOpponent(user.Username).Username)
	}
	// Send a message to the user about the Pokemon change
	um.sendMessageToUser(user, fmt.Sprintf("%s's %s has been knocked out. %s's new active Pokemon is %s.", user.Username, user.Pokemons[knockedOutIndex], user.Username, user.Pokemons[nextPokemonIndex]))

}

func (um *UserManager) switchTurn() {
	// Switch the current turn to the other player
	playerNames := um.getPlayerNames()
	for i, name := range playerNames {
		if name == um.CurrentTurn {
			um.CurrentTurn = playerNames[(i+1)%len(playerNames)]
			return
		}
	}
}

func (um *UserManager) hasLost(user *User) bool {
	if user.ActiveHP == 0 && user.ActivePokemon == user.PokemonData[2] {
		return true
	}
	return false
}

func (um *UserManager) announceWinner(winner string) {
	// Announce the winner of the battle to both players
	for _, user := range um.Users {
		_, err := user.Conn.Write([]byte(fmt.Sprintf("The winner is %s!", winner)))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", user.Username, err)
		}
	}
}

func (um *UserManager) UpdatePokemonData(username, pokemonName string, pokemonIndex int) error {
	user, exists := um.Users[username]
	if !exists {
		return fmt.Errorf("user %s not found", username)
	}
	pokemonID := utils.PokeMap[pokemonName]

	// Construct the path to the Pokemon data file relative to the current file
	//pokemonDataFilePath := filepath.Join("..", "internal", "models", "monsters", "data", fmt.Sprintf("%d.json", pokemonID))
	pokemonDataFilePath := fmt.Sprintf("../internal/models/monsters/data/%s.json", pokemonID)
	pokemonData, err := readPokemonJSONData(pokemonDataFilePath)
	if err != nil {
		return fmt.Errorf("error reading Pokemon data: %v", err)
	}

	// Update the user's Pokemon information
	if len(user.PokemonData) < pokemonIndex {
		user.PokemonData = append(user.PokemonData, pokemonData)
	} else {
		user.PokemonData[pokemonIndex-1] = pokemonData
	}

	return nil
}

//	func getPokemonIDFromName(pokemonName string) int {
//		// Construct the path to the pokemonNames.json file relative to the main.go file
//		jsonFilePath := filepath.Join("..", "internal", "models", "pokemonNames.json")
//
//		// Read the pokemonNames.json file
//		data, err := os.ReadFile(jsonFilePath)
//		if err != nil {
//			fmt.Printf("Error reading pokemonNames.json file: %v\n", err)
//			return 0
//		}
//
//		// Unmarshal the JSON data into a slice of strings
//		var pokemonNames []string
//		err = json.Unmarshal(data, &pokemonNames)
//		if err != nil {
//			fmt.Printf("Error unmarshaling pokemonNames.json data: %v\n", err)
//			return 0
//		}
//
//		// Search for the Pokemon name in the slice and return the index (which is the ID)
//		for i, name := range pokemonNames {
//			if name == pokemonName {
//				return i + 1 // The IDs start from 1, not 0
//			}
//		}
//
//		fmt.Printf("Pokemon name '%s' not found in pokemonNames.json\n", pokemonName)
//		return 0
//	}
func readPokemonJSONData(filePath string) (*PokemonData, error) {
	// Read the JSON data from the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pokemonData PokemonData
	err = json.Unmarshal(data, &pokemonData)
	if err != nil {
		return nil, err
	}

	return &pokemonData, nil
}
