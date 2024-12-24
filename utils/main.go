package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	// Path to the folder containing the JSON files
	folderPath := "../internal/models/skim_monsters/data/"

	// Initialize map
	fileContents := make(map[string]string)

	// Loop through each file from 1.json to 649.json
	for i := 1; i <= 649; i++ {
		// Construct filename
		filename := fmt.Sprintf("%d.json", i)
		filePath := folderPath + filename

		// Read file content
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			continue
		}

		// Convert content to string

		var data map[string]interface{}
		// Unmarshal JSON content
		err = json.Unmarshal(content, &data)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON from file %s: %v\n", filename, err)
			continue
		}

		// Add content to map
		fileContents[filename] = data["name"].(string)
	}

	// Print map in Go map initialization format
	fmt.Println("map[string]string{")
	for filename, content := range fileContents {
		fmt.Printf("\"%s\": \"%s\",\n", content, filename[:len(filename)-5])
	}
	fmt.Println("}")
}
