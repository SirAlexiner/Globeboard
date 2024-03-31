package main

import (
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/library"
	"globeboard/internal/handlers/endpoint/statistics"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get the port from the environment variable or set default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Define HTTP endpoints
	http.HandleFunc(Paths.Root, handlers.EmptyHandler)
	http.HandleFunc(Endpoints.Library, library.LibHandler)
	http.HandleFunc(Endpoints.ApiKey, library.APIKeyHandler)
	http.HandleFunc(Endpoints.BookCount, statistics.BookCountHandler)
	http.HandleFunc(Endpoints.Readership, statistics.ReadershipHandler)
	http.HandleFunc(Endpoints.Status, statistics.StatusHandler)
	http.HandleFunc(Endpoints.SupportedLanguages, statistics.SupportedLanguagesHandler)

	// Start the HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
