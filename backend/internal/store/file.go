package store

import (
	"encoding/json"
	"os"

	"capuchin/internal/models" // Import our new models package
)

// Load reads the file and returns the list of todos
func Load(filePath string) ([]models.Todo, error) {
	var todos []models.Todo

	data, err := os.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, generic error, return empty list
		return []models.Todo{}, nil
	}

	err = json.Unmarshal(data, &todos)
	return todos, err
}

// Save writes the list of todos to the file
func Save(filePath string, todos []models.Todo) error {
	data, _ := json.MarshalIndent(todos, "", "  ")
	return os.WriteFile(filePath, data, 0644)
}
