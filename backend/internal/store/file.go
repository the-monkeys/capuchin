package store

import (
	"encoding/json"
	"os"

	"capuchin/internal/models" // Import our models package
)

// Load reads the file and returns the list of todos
func Load(filePath string) ([]models.Todo, error) {

	// Define a list of todos
	var todos []models.Todo

	// Read the db file
	data, err := os.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, generic error, return empty list
		return []models.Todo{}, nil
	}

	// Unmarshal the data
	err = json.Unmarshal(data, &todos)
	return todos, err
}

// Save writes the list of todos to the file
func Save(filePath string, todos []models.Todo) error {
	data, _ := json.MarshalIndent(todos, "", "  ")
	return os.WriteFile(filePath, data, 0644)
}
