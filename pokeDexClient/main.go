package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8123")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// Get a Pokemon name from the user
	var pokemonName string
	fmt.Print("Enter a Pokemon name: ")
	fmt.Scanln(&pokemonName)

	// Send the Pokemon name to the server
	fmt.Fprintf(conn, pokemonName+"\n")

	// Read the entire response from the server
	response, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal("Error reading from server:", err)
	}

	// Print the response
	fmt.Println(string(response))
}
