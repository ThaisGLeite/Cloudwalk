package report

import (
	"encoding/json"
	"log"
	"os"
	"quake/models"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	games := []*models.Game{
		{
			GameNumber:   1,
			TotalKills:   0,
			Players:      map[string]*models.Player{},
			KillsByMeans: map[string]int{},
		},
		{
			GameNumber: 2,
			TotalKills: 11,
			Players: map[string]*models.Player{
				"Isgalamido": {Name: "Isgalamido", Kills: 3},
				"Mocinha":    {Name: "Mocinha", Kills: 0},
			},
			KillsByMeans: map[string]int{
				"MOD_FALLING":       1,
				"MOD_ROCKET_SPLASH": 3,
				"MOD_TRIGGER_HURT":  7,
			},
		},
		{
			GameNumber:   3,
			TotalKills:   0,
			Players:      map[string]*models.Player{},
			KillsByMeans: map[string]int{},
		},
	}

	// Create a temporary file for testing
	tmpFile := "games_test.json"

	// Write the games to the JSON file
	err := WriteJSON(games, tmpFile)
	if err != nil {
		t.Fatalf("Failed to write games to JSON: %v", err)
	}
	defer func() {
		// Remove the temporary file after the test
		err := os.Remove(tmpFile)
		if err != nil {
			log.Printf("Failed to remove temporary file: %v", err)
		}
	}()

	// Read the contents of the JSON file
	file, err := os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	fileSize := fileInfo.Size()
	fileData := make([]byte, fileSize)

	_, err = file.Read(fileData)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data
	var readGames []*models.Game
	err = json.Unmarshal(fileData, &readGames)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Compare the original games and read games
	if len(readGames) != len(games) {
		t.Errorf("Number of games mismatch. Expected %d, got %d", len(games), len(readGames))
	}

	for i := 0; i < len(games); i++ {
		if !compareGames(games[i], readGames[i]) {
			t.Errorf("Game mismatch at index %d", i)
		}
	}
}

// compareGames compares two games for equality
func compareGames(g1, g2 *models.Game) bool {
	if g1.GameNumber != g2.GameNumber {
		return false
	}

	if g1.TotalKills != g2.TotalKills {
		return false
	}

	if len(g1.Players) != len(g2.Players) {
		return false
	}

	for k, v1 := range g1.Players {
		v2, ok := g2.Players[k]
		if !ok {
			return false
		}

		if v1.Name != v2.Name || v1.Kills != v2.Kills {
			return false
		}
	}

	if len(g1.KillsByMeans) != len(g2.KillsByMeans) {
		return false
	}

	for k, v1 := range g1.KillsByMeans {
		v2, ok := g2.KillsByMeans[k]
		if !ok {
			return false
		}

		if v1 != v2 {
			return false
		}
	}

	return true
}
