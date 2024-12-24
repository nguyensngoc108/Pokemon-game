package _type

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nguyensngoc108/pokemon-game/internal/models"
)

// Move represents a Pokémon move.

// InputData represents the structure of the input text file.
type InputData struct {
	Docs []models.Types `json:"docs"`
	Seq  int            `json:"seq"`
}

func crawl() {
	// URL to fetch the data from
	url := fmt.Sprintf("https://pokedex.org/assets/types.txt")

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %s", err)
	}

	// Set the headers
	req.Header.Set("Referer", "https://pokedex.org/js/worker.js")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s", err)
	}

	// Split the content into individual JSON strings
	parts := strings.Split(string(content), "\n")

	// Initialize a slice to hold all moves

	// Iterate over each part and unmarshal the JSON into InputData
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}

		var inputData InputData
		err := json.Unmarshal([]byte(part), &inputData)
		if err != nil {
			log.Printf("Failed to unmarshal part: %s\nError: %s", part, err)
			continue
		}

		// Save each move to a separate JSON file
		for _, move := range inputData.Docs {
			filePath := fmt.Sprintf("./api/v1/type/%s/", move.ID)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				// Folder does not exist, create it
				err := os.MkdirAll(filePath, 0755) // 0755 is the permission mode for the created folder
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Folder created: %s", filePath)
			} else if err != nil {
				// Some error occurred while checking folder existence
				log.Fatal(err)
			} else {
				log.Printf("Folder already exists: %s", filePath)
			}
			filename := fmt.Sprintf("%s/poke.json", filePath)
			moveJSON, err := json.MarshalIndent(move, "", "  ")
			if err != nil {
				log.Printf("Failed to marshal move to JSON: %s\nError: %s", move.ID, err)
				continue
			}

			err = ioutil.WriteFile(filename, moveJSON, 0644)
			if err != nil {
				log.Printf("Failed to write move to file: %s\nError: %s", filename, err)
				continue
			}
		}
	}

	fmt.Println("Moves have been saved to individual JSON files.")

}
