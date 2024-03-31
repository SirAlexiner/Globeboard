// Package handlers provides HTTP request handlers for routing and handling requests within the Gutendex API.
package handlers

import (
	"fmt"
	"globeboard/internal/utils/constants/Endpoints"
	"net/http"
)

// EmptyHandler handles every request to the root path by redirecting to the endpoints.
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		emptyMethod(w)
	}
}

// emptyMethod handles the response for any methods sent to the root path.
func emptyMethod(w http.ResponseWriter) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Setting header status to dummy status 418 I'm a Teapot
	// This is fun, but also a good indicator for the web service that the API is running
	w.WriteHeader(http.StatusTeapot)

	anchorStart := "<br><a href=\""
	anchorEnd := "</a>\n"

	// Offer information for redirection to endpoints
	output := "This service does not provide any functionality on root level.\nPlease use endpoints:\n" +
		anchorStart + Endpoints.Library + "\">" + Endpoints.Library + anchorEnd +
		anchorStart + Endpoints.BookCount + "\">" + Endpoints.BookCount + anchorEnd +
		anchorStart + Endpoints.Readership + "\">" + Endpoints.Readership + anchorEnd +
		anchorStart + Endpoints.Status + "\">" + Endpoints.Status + anchorEnd +
		anchorStart + Endpoints.SupportedLanguages + "\">" + Endpoints.SupportedLanguages + anchorEnd

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
