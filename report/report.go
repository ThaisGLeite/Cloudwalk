package report

import (
	"encoding/json"
	"fmt"
	"os"
	"quake/logger"
	"quake/models"
)

// WriteJSON writes the games data to a JSON file
func WriteJSON(games []*models.Game, filePath string) error {
	// Create a new JSON file
	logger.Log.Info("Creating JSON file")
	jsonFile, err := os.Create(filePath)
	if err != nil {
		logger.Log.Error("Failed to create JSON file", err)
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer jsonFile.Close()

	// Create a new JSON encoder that writes to the file
	logger.Log.Info("Encoding games to JSON")
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "\t")

	// Write the games to the file in JSON format
	if err := encoder.Encode(games); err != nil {
		logger.Log.Error("Failed to encode games to JSON", err)
		return fmt.Errorf("failed to encode games to JSON: %w", err)
	}

	logger.Log.Info("Finished writing JSON file")
	return nil
}
