package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

type Pokemon struct {
    Name              string   `json:"name"`
    HP                int      `json:"hp"`
    Attack            int      `json:"attack"`
    Defense           int      `json:"defense"`
    SpecialAttack     int      `json:"sp_atk"`
    SpecialDefense    int      `json:"sp_def"`
    Speed             int      `json:"speed"`
    ElementalMultiplier float64 `json:"elemental_multiplier"`
    Level             int      `json:"level"`
    Abilities         []string `json:"abilities"`
    Types             []string `json:"types"`
}

type Player struct {
    Username string
    Conn     net.Conn
    Pokemons []Pokemon
    ActivePokemon int
}

type GameServer struct {
    Players []*Player
    PokemonData map[string]Pokemon
    mu      sync.Mutex
    turn    int
}

func NewGameServer(pokemonData map[string]Pokemon) *GameServer {
    return &GameServer{
        Players: make([]*Player, 0, 2),
        PokemonData: pokemonData,
        turn: 0,
    }
}

func (gs *GameServer) AddPlayer(username string, conn net.Conn) *Player {
    gs.mu.Lock()
    defer gs.mu.Unlock()

    player := &Player{
        Username: username,
        Conn:     conn,
        Pokemons: make([]Pokemon, 0, 3),
    }
    gs.Players = append(gs.Players, player)
    return player
}

func (gs *GameServer) StartBattle() {
    gs.mu.Lock()
    defer gs.mu.Unlock()

    if len(gs.Players) < 2 {
        return
    }

    player1 := gs.Players[0]
    player2 := gs.Players[1]

    fmt.Fprintf(player1.Conn, "Battle started! You are fighting against %s\n", player2.Username)
    fmt.Fprintf(player2.Conn, "Battle started! You are fighting against %s\n", player1.Username)

    // Determine who goes first based on speed
    if player1.Pokemons[player1.ActivePokemon].Speed > player2.Pokemons[player2.ActivePokemon].Speed {
        fmt.Fprintf(player1.Conn, "You go first!\n")
        fmt.Fprintf(player2.Conn, "Opponent goes first!\n")
        gs.turn = 0
    } else {
        fmt.Fprintf(player1.Conn, "Opponent goes first!\n")
        fmt.Fprintf(player2.Conn, "You go first!\n")
        gs.turn = 1
    }
}

func (gs *GameServer) HandleConnection(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            log.Println("Error reading from connection:", err)
            return
        }

        data := string(buf[:n])
        parts := strings.Split(data, " ")
        if len(parts) < 2 {
            continue
        }

        username := parts[0]
        action := parts[1]

        switch action {
        case "join":
            player := gs.AddPlayer(username, conn)
            fmt.Fprintf(conn, "Player %s joined the game\n", player.Username)
            if len(gs.Players) == 2 {
                gs.StartBattle()
            }
        case "pokemon":
            if len(parts) < 3 {
                continue
            }
            player := gs.GetPlayerByUsername(username)
            if player == nil {
                continue
            }
            pokemonNumber := len(player.Pokemons)
            if pokemonNumber >= 3 {
                continue
            }
            pokemonName := parts[2]
            pokemon, exists := gs.PokemonData[pokemonName]
            if !exists {
                fmt.Fprintf(conn, "Pokemon %s not found\n", pokemonName)
                continue
            }
            pokemon.Level = rand.Intn(100) + 1
            pokemon.ElementalMultiplier = 0.5 + rand.Float64()*0.5
            player.Pokemons = append(player.Pokemons, pokemon)
            fmt.Fprintf(conn, "Pokemon %s added to your team\n", pokemonName)
        case "move":
            if len(parts) < 3 {
                continue
            }
            player := gs.GetPlayerByUsername(username)
            if player == nil {
                continue
            }
            if gs.Players[gs.turn].Username != username {
                fmt.Fprintf(conn, "It's not your turn!\n")
                continue
            }
            moveType := parts[2]
            gs.PerformMove(player, moveType)
            gs.turn = (gs.turn + 1) % 2
        }
    }
}

func (gs *GameServer) GetPlayerByUsername(username string) *Player {
    gs.mu.Lock()
    defer gs.mu.Unlock()

    for _, player := range gs.Players {
        if player.Username == username {
            return player
        }
    }
    return nil
}

func (gs *GameServer) PerformMove(player *Player, moveType string) {
    opponent := gs.GetOpponent(player)
    if opponent == nil {
        return
    }

    playerPokemon := player.Pokemons[player.ActivePokemon]
    opponentPokemon := opponent.Pokemons[opponent.ActivePokemon]

    var damage int
    if moveType == "normal" {
        damage = playerPokemon.Attack - opponentPokemon.Defense
    } else {
        damage = int(float64(playerPokemon.SpecialAttack) * playerPokemon.ElementalMultiplier) - opponentPokemon.SpecialDefense
    }

    if damage < 0 {
        damage = 0
    }

    opponentPokemon.HP -= damage
    if opponentPokemon.HP <= 0 {
        opponentPokemon.HP = 0
        fmt.Fprintf(opponent.Conn, "Your %s fainted!\n", opponentPokemon.Name)
        fmt.Fprintf(player.Conn, "You defeated %s's %s!\n", opponent.Username, opponentPokemon.Name)
        opponent.ActivePokemon++
        if opponent.ActivePokemon >= len(opponent.Pokemons) {
            fmt.Fprintf(player.Conn, "You won the battle!\n")
            fmt.Fprintf(opponent.Conn, "You lost the battle!\n")
            gs.EndBattle(player, opponent)
        } else {
            fmt.Fprintf(opponent.Conn, "Send out your next Pokemon!\n")
        }
    } else {
        fmt.Fprintf(opponent.Conn, "Your %s took %d damage!\n", opponentPokemon.Name, damage)
        fmt.Fprintf(player.Conn, "You dealt %d damage to %s's %s!\n", damage, opponent.Username, opponentPokemon.Name)
    }
}

func (gs *GameServer) GetOpponent(player *Player) *Player {
    gs.mu.Lock()
    defer gs.mu.Unlock()

    for _, p := range gs.Players {
        if p != player {
            return p
        }
    }
    return nil
}

func (gs *GameServer) EndBattle(winner, loser *Player) {
    // Calculate experience points
    totalExp := 0
    for _, pokemon := range loser.Pokemons {
        totalExp += pokemon.Level * 10
    }
    expPerPokemon := totalExp / len(winner.Pokemons)
    for i := range winner.Pokemons {
        winner.Pokemons[i].Level += expPerPokemon
    }

    // Reset game state
    gs.Players = make([]*Player, 0, 2)
}

func loadPokemonData(path string) (map[string]Pokemon, error) {
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, err
    }

    pokemonData := make(map[string]Pokemon)
    for _, file := range files {
        if filepath.Ext(file.Name()) == ".json" {
            data, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
            if err != nil {
                return nil, err
            }

            var pokemon struct {
                Monster Pokemon `json:"monster"`
            }
            if err := json.Unmarshal(data, &pokemon); err != nil {
                return nil, err
            }

            pokemonData[pokemon.Monster.Name] = pokemon.Monster
        }
    }

    return pokemonData, nil
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // Load Pokémon data
    pokemonData, err := loadPokemonData("../internal/models/monsters/data")
    if err != nil {
        log.Fatal("Error loading Pokémon data:", err)
    }

    gameServer := NewGameServer(pokemonData)

    listener, err := net.Listen("tcp", ":8000")
    if err != nil {
        log.Fatal("Error starting server:", err)
    }
    defer listener.Close()

    fmt.Println("Server started. Listening on port 8000...")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Failed to accept connection:", err)
            continue
        }
        go gameServer.HandleConnection(conn)
    }
}