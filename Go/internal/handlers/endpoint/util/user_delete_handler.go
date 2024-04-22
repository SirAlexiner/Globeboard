// Package util provides HTTP handlers for user and API key management within the application.
package util

import (
	"context"
	authenticate "globeboard/auth"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"log"
	"net/http"
)

// UserDeletionHandler handles HTTP requests for user deletion.
func UserDeletionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteUser(w, r) // Handle DELETE requests
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.UserDeletionID, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is:\n"+http.MethodDelete, http.StatusNotImplemented)
		return
	}
}

// deleteUser processes the user deletion using the user ID from the request path.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")    // Extract user ID from the URL path.
	if ID == "" || ID == " " { // Check if the user ID is provided.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.UserDeletionID)
		http.Error(w, "Please provide User ID", http.StatusBadRequest) // Return an error if user ID is missing.
		return
	}

	ctx := context.Background() // Create a new background context.

	err := authenticate.Client.DeleteUser(ctx, ID) // Attempt to delete user in Firebase.
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)               // Log the error.
		http.Error(w, err.Error(), http.StatusInternalServerError) // Report deletion error.
		return
	}

	w.WriteHeader(http.StatusNoContent)               // Set HTTP status to 204 No Content on successful deletion.
	log.Printf("Successfully deleted user: %v\n", ID) // Log successful deletion.
}
