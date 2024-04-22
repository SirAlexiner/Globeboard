// Package handlers provide HTTP request handlers for routing and handling requests within the application.
package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	ISE = "Internal Server Error" // ISE defines the error message returned when an internal server error occurs.
)

// EmptyHandler serves the "root.html" file at the root URL ("/") and returns a 404 Not Found error for other paths.
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r) // Return a 404 Not Found error if the path is not the root.
		return
	}

	filePath := "./web/root.html" // Path to the "root.html" file.

	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Set the "Content-Type" header.
	w.WriteHeader(http.StatusSeeOther)                         // Set the status code to 303 See Other,
	// indicating redirection.

	file, err := os.Open(filePath) // Open the "root.html" file.
	if err != nil {
		log.Printf("%s: Error opening root file: %v", r.RemoteAddr, err) // Log error if file opening fails.
		http.Error(w, ISE, http.StatusInternalServerError)               // Return a 500 Internal Server Error if the file cannot be opened.
		return
	}
	defer func(file *os.File) { // Ensure the file is closed after serving it.
		err := file.Close()
		if err != nil {
			log.Printf("%s: Error closing root file: %v", r.RemoteAddr, err) // Log error if file closing fails.
			http.Error(w, ISE, http.StatusInternalServerError)               // Return a 500 Internal Server Error
			// if the file cannot be closed.
		}
	}(file)

	_, err = io.Copy(w, file) // Copy the file content to the response writer.
	if err != nil {
		log.Printf("%s: Error copying root file to ResponseWriter: %v", r.RemoteAddr, err) // Log error if copying fails.
		http.Error(w, ISE, http.StatusInternalServerError)                                 // Return a 500 Internal Server Error
		// if content cannot be copied.
		return
	}
}
