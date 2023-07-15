// main_test.go
package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"quake/handler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	h := handler.NewHandler()

	server := httptest.NewServer(http.HandlerFunc(h.Upload))
	defer server.Close()

	// Create a buffer to store our request body
	body := &bytes.Buffer{}

	// Create a multipart writer to write the file to the buffer
	writer := multipart.NewWriter(body)

	// Open the log file
	file, err := os.Open("../testdata/test.log")
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a form file writer for the log file
	formFile, err := writer.CreateFormFile("qgames", filepath.Base(file.Name()))
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		t.Fatalf("Error copying file to form file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Error closing multipart writer: %v", err)
	}

	// Create a new file upload request with our body, which contains the log file
	req, err := http.NewRequest("POST", server.URL+"/upload", body)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Set the content type of the request to be multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected response status OK")

	// TODO: Further assertions can be made here based on the response body.
}
