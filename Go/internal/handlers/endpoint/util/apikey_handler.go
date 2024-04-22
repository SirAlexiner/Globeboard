// Package util provides HTTP handlers for user and API key management within the application.
package util

import (
	"context"
	"encoding/json"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"log"
	"net/http"
)

// APIKeyHandler routes API Key management requests to the appropriate functions based on the HTTP method.
func APIKeyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // Handle GET requests
		handleApiKeyGetRequest(w, r)
	case http.MethodDelete: // Handle DELETE requests
		handleApiKeyDeleteRequest(w, r)
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.ApiKey, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint are: GET, DELETE", http.StatusNotImplemented)
		return
	}
}

// handleApiKeyDeleteRequest handles the deletion of an API key.
func handleApiKeyDeleteRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token") // Retrieve the token from query parameters.

	UUID := r.Header.Get("Authorization") // Retrieve the UUID from the Authorization header.

	ctx := context.Background() // Create a new background context.

	_, err := authenticate.Client.GetUser(ctx, UUID) // Verify the UUID with Firebase Authentication.
	if err != nil {
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.ApiKey)
		log.Printf("Error verifying UUID: %v\n", err)
		http.Error(w, "Not Authorized", http.StatusUnauthorized) // Respond with unauthorized if UUID is invalid.
		return
	}

	if token == "" || token == " " { // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.ApiKey)
		http.Error(w, "Please specify API Key to delete: '?token={API_Key}'", http.StatusBadRequest)
		return
	}

	err = db.DeleteApiKey(r.RemoteAddr, UUID, token) // Attempt to delete the API key.
	if err != nil {
		log.Printf("%s: Error deleting API Key: %v", r.RemoteAddr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Respond with internal server error if deletion fails.
		return
	}

	w.WriteHeader(http.StatusNoContent) // Respond with no content on successful deletion.
}

// handleApiKeyGetRequest handles the creation and retrieval of a new API key.
func handleApiKeyGetRequest(w http.ResponseWriter, r *http.Request) {
	UDID := _func.GenerateUID(constants.DocIdLength)    // Generate a unique document ID.
	key := _func.GenerateAPIKey(constants.ApiKeyLength) // Generate a new API key.

	UUID := r.Header.Get("Authorization") // Retrieve the UUID from the Authorization header.

	w.Header().Set("Content-Type", "application/json") // Set the content type of the response to application/json.

	ctx := context.Background() // Create a new background context.

	_, err := authenticate.Client.GetUser(ctx, UUID) // Verify the UUID with Firebase Authentication.
	if err != nil {
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.ApiKey)
		log.Printf("%s: Error verifying UUID: %v\n", r.RemoteAddr, err)
		http.Error(w, "Not Authorized", http.StatusUnauthorized) // Respond with unauthorized if UUID is invalid.
		return
	}

	err = db.AddApiKey(r.RemoteAddr, UDID, UUID, key) // Attempt to add the new API key to the database.
	if err != nil {
		log.Printf("%s: Error creating API Key: %v", r.RemoteAddr, err)
		errorMessage := fmt.Sprintf("Error creating API Key: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError) // Respond with internal server error if addition fails.
		return
	}

	w.WriteHeader(http.StatusCreated) // Set HTTP status to 201 Created on successful API key creation.

	response := map[string]string{
		"token": key, // Include the new API key in the response.
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError) // Handle errors in JSON encoding.
		return
	}
}
