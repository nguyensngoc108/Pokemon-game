package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net"
    "os"
    "strings"
    "sync"
    "time"
)

const (
    worldSize       = 1000
    maxPokemons     = 200
    spawnInterval   = time.Minute
    despawnInterval = 5 * time.Minute
)

type Position struct {
    X, Y int
}

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
    Password string
    Conn     net.Conn
    Position Position
    Pokemons []Pokemon
}

type GameWorld struct {
    Players  map[string]*Player
    Pokemons map[Position]Pokemon
    mu       sync.Mutex
}

func NewGameWorld() *GameWorld {
    return &GameWorld{
        Players:  make(map[string]*Player),
        Pokemons: make(map[Position]Pokemon),
    }
}

func (gw *GameWorld) AddPlayer(username, password string, conn net.Conn) *Player {
    gw.mu.Lock()
    defer gw.mu.Unlock()

    player := &Player{
        Username: username,
        Password: password,
        Conn:     conn,
        Position: Position{
            X: rand.Intn(worldSize),
            Y: rand.Intn(worldSize),
        },
        Pokemons: []Pokemon{},
    }
    gw.Players[username] = player
    return player
}

func (gw *GameWorld) MovePlayer(username, direction string) {
    gw.mu.Lock()
    defer gw.mu.Unlock()

    player, exists := gw.Players[username]
    if !exists {
        fmt.Println("Player not found:", username)
        return
    }

    switch direction {
    case "up":
        if player.Position.Y > 0 {
            player.Position.Y--
        }
    case "down":
        if player.Position.Y < worldSize-1 {
            player.Position.Y++
        }
    case "left":
        if player.Position.X > 0 {
            player.Position.X--
        }
    case "right":
        if player.Position.X < worldSize-1 {
            player.Position.X++
        }
    }

    // Check for Pokémon capture
    if pokemon, exists := gw.Pokemons[player.Position]; exists {
        if len(player.Pokemons) < maxPokemons {
            player.Pokemons = append(player.Pokemons, pokemon)
            delete(gw.Pokemons, player.Position)
            fmt.Printf("%s captured a %s\n", username, pokemon.Name)
            savePlayerData(player)
        }
    }
}

func (gw *GameWorld) SpawnPokemons() {
    for {
        time.Sleep(spawnInterval)
        gw.mu.Lock()
        for i := 0; i < 50; i++ {
            pos := Position{
                X: rand.Intn(worldSize),
                Y: rand.Intn(worldSize),
            }
            pokemonName := getRandomPokemonName()
            pokemon, err := loadPokemonData(pokemonName)
            if err != nil {
                fmt.Println("Error loading Pokémon data:", err)
                continue
            }
            pokemon.ElementalMultiplier = 0.5 + rand.Float64()*0.5
            pokemon.Level = rand.Intn(100) + 1
            gw.Pokemons[pos] = pokemon
        }
        gw.mu.Unlock()
    }
}

func (gw *GameWorld) DespawnPokemons() {
    for {
        time.Sleep(despawnInterval)
        gw.mu.Lock()
        for pos := range gw.Pokemons {
            delete(gw.Pokemons, pos)
        }
        gw.mu.Unlock()
    }
}

func getRandomPokemonName() string {
    pokemonNames := loadPokemonNames()
    return pokemonNames[rand.Intn(len(pokemonNames))]
}

func loadPokemonNames() []string {
    var pokemonNames []string
    data, err := os.ReadFile("../internal/models/pokemonNames.json")
    if err != nil {
        log.Fatal("Error loading Pokémon names:", err)
    }
    err = json.Unmarshal(data, &pokemonNames)
    if err != nil {
        log.Fatal("Error unmarshalling Pokémon names:", err)
    }
    return pokemonNames
}

func loadPokemonData(name string) (Pokemon, error) {
    var pokemon Pokemon
    filePath := fmt.Sprintf("../internal/models/monsters/data/%s.json", name)
    data, err := os.ReadFile(filePath)
    if err != nil {
        return pokemon, err
    }
    err = json.Unmarshal(data, &pokemon)
    if err != nil {
        return pokemon, err
    }
    return pokemon, nil
}

func authenticatePlayer(username, password string) bool {
    var players []Player
    data, err := os.ReadFile("players.json")
    if err != nil {
        log.Fatal("Error loading players data:", err)
    }
    err = json.Unmarshal(data, &players)
    if err != nil {
        log.Fatal("Error unmarshalling players data:", err)
    }
    for _, player := range players {
        if player.Username == username && player.Password == password {
            return true
        }
    }
    return false
}

func savePlayerData(player *Player) {
    var players []Player
    data, err := os.ReadFile("players.json")
    if err != nil {
        log.Fatal("Error loading players data:", err)
    }
    err = json.Unmarshal(data, &players)
    if err != nil {
        log.Fatal("Error unmarshalling players data:", err)
    }
    for i, p := range players {
        if p.Username == player.Username {
            players[i] = *player
            break
        }
    }
    data, err = json.MarshalIndent(players, "", "  ")
    if err != nil {
        log.Fatal("Error marshalling players data:", err)
    }
    err = os.WriteFile("players.json", data, 0644)
    if err != nil {
        log.Fatal("Error saving players data:", err)
    }
}

func startTCPServer(gw *GameWorld) {
    listener, err := net.Listen("tcp", ":8000")
    if err != nil {
        fmt.Println("Failed to listen on port 8000:", err)
        return
    }
    defer listener.Close()
    fmt.Println("TCP server listening on :8000")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Failed to accept connection:", err)
            continue
        }
        go handleConnection(conn, gw)
    }
}

func handleConnection(conn net.Conn, gw *GameWorld) {
    defer conn.Close()
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading from connection:", err)
            return
        }

        data := string(buf[:n])
        parts := strings.Split(data, " ")
        if len(parts) < 3 {
            continue
        }

        username := parts[0]
        password := parts[1]
        action := parts[2]

        // Log the received username and password
        fmt.Printf("Received username: %s, password: %s\n", username, password)

        if !authenticatePlayer(username, password) {
            fmt.Fprintf(conn, "Authentication failed for user %s\n", username)
            return
        }

        fmt.Fprintf(conn, "Login successful for user %s\n", username)

        switch action {
        case "join":
            player := gw.Players[username]
            if player == nil {
                player = gw.AddPlayer(username, password, conn)
                fmt.Fprintf(conn, "Player %s joined at position (%d, %d)\n", player.Username, player.Position.X, player.Position.Y)
            } else {
                fmt.Fprintf(conn, "Player %s is already in the game\n", player.Username)
            }
        case "move":
            if len(parts) < 4 {
                continue
            }
            direction := parts[3]
            gw.MovePlayer(username, direction)
            player := gw.Players[username]
            if player == nil {
                fmt.Println("Player not found:", username)
                return
            }
            fmt.Fprintf(conn, "Player %s moved to position (%d, %d)\n", player.Username, player.Position.X, player.Position.Y)
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    gameWorld := NewGameWorld()

    go gameWorld.SpawnPokemons()
    go gameWorld.DespawnPokemons()

    startTCPServer(gameWorld)
}