// Handler package handles HTTP requests and responses.
package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"quake/logger"
	"quake/models"
	"quake/parser"
	"strings"
	"time"
)

// Handler struct holds a pointer to a parser.
type Handler struct {
	parser *parser.Parser
}

// NewHandler returns a pointer to a Handler struct.
func NewHandler() *Handler {
	return &Handler{
		parser: parser.NewParser(),
	}
}

// Upload is an HTTP handler function that processes the uploaded log file and returns parsed game data.
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// Create a new context with a deadline for the request.
	_, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	logger.Log.Info("Received a file upload request")

	// Parse the form data in the request.
	err := r.ParseMultipartForm(10 << 20) // 10 MB of maximum file size
	if err != nil {
		http.Error(w, "Error Parsing the Form", http.StatusInternalServerError)
		logger.Log.Errorf("Error parsing form: %v", err)
		return
	}

	// Get the uploaded file from the request.
	file, header, err := r.FormFile("qgames")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusInternalServerError)
		logger.Log.Errorf("Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	// Validate the file type and/or extension.
	buffer := make([]byte, 512) // 512 bytes should be enough for most formats
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		http.Error(w, "Error Reading the File", http.StatusInternalServerError)
		logger.Log.Errorf("Error reading the file: %v", err)
		return
	}
	mimetype := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimetype, "text/plain") || filepath.Ext(header.Filename) != ".log" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		logger.Log.Error("Invalid file type uploaded")
		return
	}

	// Create a temporary file to store the uploaded file.
	tempFile, err := os.CreateTemp("", "upload-*.log")
	if err != nil {
		http.Error(w, "Error Creating the Temporary File", http.StatusInternalServerError)
		logger.Log.Errorf("Error creating temporary file: %v", err)
		return
	}
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file.
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error Writing the Temporary File", http.StatusInternalServerError)
		logger.Log.Errorf("Error writing to temporary file: %v", err)
		return
	}

	// Parse the temporary log file.
	gameMap, err := h.parser.ParseLog(tempFile.Name())
	if err != nil {
		http.Error(w, "Error Parsing Log", http.StatusInternalServerError)
		logger.Log.Errorf("Error parsing log: %v", err)
		return
	}

	// Convert the map of games to a slice of games.
	games := make([]*models.Game, 0, len(gameMap))
	games = append(games, gameMap...)

	// Marshal the slice of games to JSON.
	js, err := json.Marshal(games)
	if err != nil {
		http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
		logger.Log.Errorf("Error marshalling JSON: %v", err)
		return
	}

	// Set the response header and write the JSON to the response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	logger.Log.Info("Finished writing report")
}
