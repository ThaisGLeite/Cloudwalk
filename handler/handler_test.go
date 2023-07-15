package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type TestCase struct {
	filename    string
	contentType string
	expected    int
}

func TestUploadHandler(t *testing.T) {
	testCases := []TestCase{
		{
			filename:    "../testdata/test.log",
			contentType: "text/plain",
			expected:    http.StatusOK,
		},
		{
			filename:    "../testdata/test.png",
			contentType: "image/png",
			expected:    http.StatusBadRequest,
		},
		// add more test cases if needed
	}

	for _, tc := range testCases {
		// Open the file
		file, err := os.Open(tc.filename)
		if err != nil {
			t.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		// Create a multipart writer
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("qgames", filepath.Base(file.Name()))
		if err != nil {
			t.Fatalf("Error creating form file: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("Error copying file to part: %v", err)
		}
		writer.Close()

		// Create a request
		req, err := http.NewRequest("POST", "/upload", body)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())

		// Create a ResponseRecorder
		rr := httptest.NewRecorder()

		h := NewHandler()
		handler := http.HandlerFunc(h.Upload)

		// Send the request to the handler
		handler.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != tc.expected {
			t.Errorf("Expected response status %v, got %v for file %s", tc.expected, status, tc.filename)
		}
	}
}
