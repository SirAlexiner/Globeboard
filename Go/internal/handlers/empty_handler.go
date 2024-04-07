// Package handlers provide HTTP request handlers for routing and handling requests within the Gutendex API.
package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	ISE = "Internal Server Error"
)

// EmptyHandler handles every request to the root path by redirecting to the endpoints.
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	filePath := "./web/root.html"

	// Set the "Content-Type" header.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set the status code to indicate the redirection.
	w.WriteHeader(http.StatusSeeOther)

	// Open the file.
	file, err := os.Open(filePath)
	if err != nil {
		log.Print("Error opening root file: ", err)
		http.Error(w, ISE, http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print("Error closing root file: ", err)
			http.Error(w, ISE, http.StatusInternalServerError)
			return
		}
	}(file)

	_, err = io.Copy(w, file)
	if err != nil {
		log.Print("Error copying root file to ResponseWriter: ", err)
		http.Error(w, ISE, http.StatusInternalServerError)
		return
	}
}
