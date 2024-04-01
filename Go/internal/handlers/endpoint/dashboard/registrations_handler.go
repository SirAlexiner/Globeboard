// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"net/http"
	"sort"
)

// RegistrationsHandler handles HTTP GET requests to retrieve supported languages.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSupportedLanguagesGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

// handleSupportedLanguagesGetRequest handles GET requests to retrieve supported languages.
func handleSupportedLanguagesGetRequest(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Please provide API Token", http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf("Error checking API key: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	if !exists {
		err := fmt.Sprintf("API key not accepted")
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	// Get the supported languages from self defined function
	languages, err := _func.GetSupportedLanguages()
	if err != nil {
		http.Error(w, "Error whilst retrieving languages: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Sort the language entries alphabetically by name
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Name < languages[j].Name
	})

	// Write content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode languages as JSON and send the response
	if err := json.NewEncoder(w).Encode(languages); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
