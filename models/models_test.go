package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPlayerJSONSerialization(t *testing.T) {
	player := &Player{
		Name:  "John",
		Kills: 5,
	}

	// Serialize the player to JSON
	jsonData, err := json.Marshal(player)
	if err != nil {
		t.Fatalf("Failed to serialize player to JSON: %v", err)
	}

	// Deserialize the JSON back to a player object
	var deserializedPlayer Player
	err = json.Unmarshal(jsonData, &deserializedPlayer)
	if err != nil {
		t.Fatalf("Failed to deserialize JSON to player: %v", err)
	}

	// Compare the original player and the deserialized player
	if !reflect.DeepEqual(player, &deserializedPlayer) {
		t.Errorf("Deserialized player does not match the original player")
	}
}

func TestGameJSONSerialization(t *testing.T) {
	game := &Game{
		GameNumber: 1,
		TotalKills: 10,
		Players: map[string]*Player{
			"Player1": {Name: "Player1", Kills: 5},
			"Player2": {Name: "Player2", Kills: 3},
		},
		KillsByMeans: map[string]int{
			"MOD_SHOTGUN":    5,
			"MOD_ROCKET":     3,
			"MOD_MACHINEGUN": 2,
		},
	}

	// Serialize the game to JSON
	jsonData, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("Failed to serialize game to JSON: %v", err)
	}

	// Deserialize the JSON back to a game object
	var deserializedGame Game
	err = json.Unmarshal(jsonData, &deserializedGame)
	if err != nil {
		t.Fatalf("Failed to deserialize JSON to game: %v", err)
	}

	// Compare the original game and the deserialized game
	if !reflect.DeepEqual(game, &deserializedGame) {
		t.Errorf("Deserialized game does not match the original game")
	}
}
