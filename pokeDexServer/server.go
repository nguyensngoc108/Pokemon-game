package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/nguyensngoc108/pokemon-game/internal/handlers"
	"github.com/nguyensngoc108/pokemon-game/internal/repositories"
)

// Server represents the file server
type Server struct{}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	var BasePath = "../internal/models"
	pokedexRepository := repositories.NewPokedexRepository(BasePath)
	pokedexHandler := handlers.NewPokeDexHandler(pokedexRepository)

	fileName, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading file name:", err)
		return
	}
	fileName = strings.TrimSpace(fileName)

	// Read the file content
	fileContent, err := pokedexHandler.GetPokemon(fileName)
	if err != nil {
		log.Println("Error reading file:", err)
		conn.Write([]byte("Error reading file\n"))
		return
	}

	// Send the file content to the pokeDexClient
	_, err = conn.Write(fileContent)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8123")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on port 8123...")

	server := &Server{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go server.HandleConnection(conn)
	}
}
