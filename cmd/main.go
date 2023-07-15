// main.go
package main

import (
	"net/http"
	"quake/handler"
	"quake/logger"
)

func main() {
	logger.Log.Info("Starting application")

	h := handler.NewHandler()

	http.HandleFunc("/upload", h.Upload)

	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	logger.Log.Info("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Log.Fatal(err)
	}
}
