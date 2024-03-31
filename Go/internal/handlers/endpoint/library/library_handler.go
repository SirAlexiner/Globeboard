// Package library provides handlers and utilities for managing the library API endpoints.
package library

import (
	"encoding/json"
	"fmt"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// LibHandler handles requests to the library API endpoint.
func LibHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleBookGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

// handleBookGetRequest handles GET requests to retrieve book data.
func handleBookGetRequest(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from the request.
	query := r.URL.Query()
	languagesString := query.Get("languages")
	search := query.Get("search")
	page := query.Get("page")

	// Create a new request to the Gutendex API.
	req, err := http.NewRequest(http.MethodGet, External.GutendexAPI, nil)
	if err != nil {
		err := fmt.Sprintf("error creating request: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
	}

	// Split languages string into individual languages.
	languages := strings.Split(languagesString, ",")
	// Validate languages.
	if err := _func.IsValidLanguages(languages, true, true); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if languagesString != "" {
		q := req.URL.Query()
		q.Add("languages", languagesString)
		req.URL.RawQuery = q.Encode()
	}

	// Add search query parameter if provided.
	if search != "" {
		q := req.URL.Query()
		q.Add("search", search)
		req.URL.RawQuery = q.Encode()
	}

	// Add page query parameter if provided.
	if page != "" {
		// Make sure page query is a number above 0
		if _, err := strconv.Atoi(page); err != nil {
			http.Error(w, "page parameter must be a number", http.StatusBadRequest)
			return
		}
		if page == "0" {
			http.Error(w, "page parameter must larger than 0", http.StatusBadRequest)
			return
		}
		q := req.URL.Query()
		q.Add("page", page)
		req.URL.RawQuery = q.Encode()
	}

	// Set content type.
	req.Header.Add("content-type", "application/json")

	// Create an HTTP client.
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request to Gutendex API.
	res, err := client.Do(req)
	if err != nil {
		err := fmt.Errorf("error in response: %s", err.Error())
		log.Println(err)
	}

	// Decode response into BookResponse struct.
	var response structs.BookResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		http.Error(w, "Failed to decode books response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Trim API URL prefix from Next and Previous fields.
	if response.Next != "" {
		result := strings.TrimPrefix(response.Next, External.GutendexAPI)
		response.Next = result
	}

	if response.Previous != "" {
		result := strings.TrimPrefix(response.Previous, External.GutendexAPI)
		response.Previous = result
	}

	// Set content type header.
	w.Header().Set("Content-Type", "application/json")

	// Encode books as JSON and send the response.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
