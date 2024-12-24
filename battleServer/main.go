package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"github.com/nguyensngoc108/pokemon-game/battleServer/usermanager"
)

func main() {
	startTCPServer()
}

func startTCPServer() {
	// Listen on TCP port
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Failed to listen on port 8000:", err)
		return
	}
	defer listener.Close()
	userManager := usermanager.NewUserManager()
	usermanager.SetUserManagerInstance(userManager)
	fmt.Println("TCP server listening on :8000")

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, userManager)
	}
}

func handleConnection(conn net.Conn, userManager *usermanager.UserManager) {
	defer conn.Close()
	// Store the userManager instance in the singleton
	usermanager.GetUserManagerInstance().Users = userManager.Users
	buf := make([]byte, 1024)
	for {
		// Read data from the connection
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// Parse the received data
		data := strings.TrimSpace(string(buf[:n]))
		parts := strings.Split(data, " ")

		switch len(parts) {
		case 1:
			// New pokeDexClient connection
			username := parts[0]
			user := userManager.AddUser(username, conn)
			fmt.Printf("New pokeDexClient connected: %s\n", user.Username)
		case 3:
			// Received Pokemon information
			username := parts[0]
			pokemonNumber, _ := strconv.Atoi(parts[1])
			pokemonName := parts[2]

			err := userManager.UpdatePokemonData(username, pokemonName, pokemonNumber)
			if err != nil {
				fmt.Printf("Error updating Pokemon data for %s: %v\n", username, err)
				continue
			}

			userManager.UpdatePokemons(username, pokemonName, pokemonNumber)
			fmt.Printf("%s added Pokemon %d: %s\n", username, pokemonNumber, pokemonName)

			if userManager.AllPokemonsProvided() && len(userManager.Users) == 2 {
				userManager.StartBattle()
			}
		case 4:
			if userManager.AllPokemonsProvided() {
				moveType := parts[1]
				username := parts[0]
				// Process the move information and send the result back to the clients
				userManager.PerformBattle(moveType, username)
			}
		}
	}
}
