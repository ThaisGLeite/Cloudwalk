package parser

import (
	"quake/models"
	"testing"
)

func TestSplitIntoChunks(t *testing.T) {
	// Modify this with the path to a test file in your system
	testFilePath := "../testdata/test.log"

	_, err := splitIntoChunks(testFilePath)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

func TestParseChunk(t *testing.T) {
	// Test input. Modify it based on your needs.
	input := `Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
	           InitGame: 
	           Kill: 1022 3 22: <world> killed Dono da Bola by MOD_TRIGGER_HURT`

	// The game number for this test
	gameNumber := 1

	expectedGame := &models.Game{
		GameNumber: gameNumber,
		Players: map[string]*models.Player{
			"Isgalamido":   {Name: "Isgalamido", Kills: -1},
			"Dono da Bola": {Name: "Dono da Bola", Kills: -1},
		},
		KillsByMeans: map[string]int{"MOD_TRIGGER_HURT": 2},
		TotalKills:   2,
	}

	game, err := parseChunk(input, gameNumber)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if game.GameNumber != expectedGame.GameNumber {
		t.Fatalf("Expected game number %v, but got: %v", expectedGame.GameNumber, game.GameNumber)
	}

	// Add more checks based on your requirements
	// ...
}

func TestParseLog(t *testing.T) {
	// Modify this with the path to a test file in your system
	testFilePath := "../testdata/test.log"

	parser := NewParser()

	_, err := parser.ParseLog(testFilePath)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}
