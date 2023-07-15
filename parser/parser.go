package parser

import (
	"bufio"
	"fmt"
	"os"
	"quake/logger"
	"quake/models"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Parser struct {
	// TODO: any state needed for parsing
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseLog(filePath string) ([]*models.Game, error) {
	// Read the log file and split it into chunks
	chunks, err := splitIntoChunks(filePath)
	if err != nil {
		logger.Log.Fatalf("Error dividing log into chunks: %v", err)
	}

	// Create a channel to collect the games
	games := make(chan *models.Game, len(chunks))

	// Create a channel to send chunks to workers
	tasks := make(chan string, len(chunks))

	// Create a channel to send game numbers to workers
	gameNumbers := make(chan int, len(chunks))

	// Start a fixed number of workers
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ { // 10 is the number of workers
		wg.Add(1)
		go func() {
			defer wg.Done()

			for chunk := range tasks {
				gameNumber := <-gameNumbers
				game, err := parseChunk(chunk, gameNumber)
				if err != nil {
					logger.Log.Printf("Error parsing chunk: %v", err)
					continue
				}
				games <- game
			}
		}()
	}

	// Send chunks and game numbers to the tasks and gameNumbers channels
	for i, chunk := range chunks {
		tasks <- chunk
		gameNumbers <- i + 1
	}
	close(tasks)
	close(gameNumbers)

	// Wait for all workers to finish, then close the games channel
	go func() {
		wg.Wait()
		close(games)
	}()

	// Collect the games into a slice
	result := []*models.Game{}
	for game := range games {
		result = append(result, game)
	}

	// Sort the games based on the game number
	sort.Slice(result, func(i, j int) bool {
		return result[i].GameNumber < result[j].GameNumber
	})

	return result, nil
}

// splitIntoChunks takes a file path as input and returns a slice of strings,
// where each string represents a chunk of the game log file. It returns an error if any occurs during this process.
func splitIntoChunks(filePath string) ([]string, error) {
	// Open the log file and ensure the file is properly closed after the function returns
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Prepare a slice to hold chunks of log data
	scanner := bufio.NewScanner(file)
	chunks := []string{}
	var currentChunk strings.Builder

	// Loop through all lines in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line indicates a new game
		if strings.Contains(line, "InitGame:") {
			// If the current chunk is not empty, add it to the chunks slice
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				// Clear the currentChunk for the new game
				currentChunk.Reset()
			}
		}

		// Add the current line to the current chunk
		currentChunk.WriteString(line)
		currentChunk.WriteString("\n")
	}

	// If the scanner encountered any error during reading the file, return the error
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %v", err)
	}

	// If the current chunk is not empty after scanning the whole file, add it to the chunks
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks, nil
}

func parseChunk(chunk string, gameNumber int) (*models.Game, error) {
	lines := strings.Split(chunk, "\n")
	game := &models.Game{
		GameNumber:   gameNumber,
		Players:      make(map[string]*models.Player),
		KillsByMeans: make(map[string]int),
	}

	killRegexp := regexp.MustCompile(`Kill: \d+ \d+ \d+: (.+) killed (.+) by (.+)`)

	for _, line := range lines {
		if strings.Contains(line, "Kill:") {
			matches := killRegexp.FindStringSubmatch(line)
			if len(matches) == 4 {
				killerName := matches[1]
				victimName := matches[2]
				meansOfDeath := matches[3]

				if _, exists := game.Players[victimName]; !exists {
					game.Players[victimName] = &models.Player{Name: victimName}
				}

				if killerName == "<world>" {
					game.Players[victimName].Kills--
				} else if killerName != victimName { // Only increment kills if the killer isn't the victim
					if _, exists := game.Players[killerName]; !exists {
						game.Players[killerName] = &models.Player{Name: killerName}
					}
					game.Players[killerName].Kills++
				}
				game.KillsByMeans[meansOfDeath]++
				game.TotalKills++
			}
		}
	}

	return game, nil
}
