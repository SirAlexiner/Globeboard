// Package dashboard provides handlers for managing dashboard-related functionalities through HTTP endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"log"
	"net/http"
)

// NotificationsIdHandler handles HTTP requests related to specific notification settings by ID.
func NotificationsIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // Handle GET request
		handleNotifGetRequest(w, r)
	case http.MethodDelete: // Handle DELETE request
		handleNotifDeleteRequest(w, r)
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.NotificationsID, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint are:\n"+http.MethodGet+"\n"+http.MethodDelete, http.StatusNotImplemented)
		return
	}
}

// handleNotifGetRequest processes GET requests to retrieve a specific notification webhook by its ID.
func handleNotifGetRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")     // Retrieve the ID from the URL path.
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Retrieve token from the URL query parameters.
	if token == "" {            // Check if a token is provided.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.NotificationsID)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token) // Retrieve UUID associated with the API token.
	if UUID == "" {                 // Check if UUID is valid.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.NotificationsID)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" || ID == " " { // Check if the ID is valid.
		log.Printf(constants.ClientConnectNoID, r.RemoteAddr, r.Method, Endpoints.NotificationsID)
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	hook, err := db.GetSpecificWebhook(ID, UUID) // Retrieve the specific webhook by ID and UUID.
	if err != nil {
		log.Printf("%s: Error getting webhook from database: %v", r.RemoteAddr, err)
		err := fmt.Sprintf("Error getting webhook from database: %v", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON) // Set the content type of the response.

	w.WriteHeader(http.StatusOK) // Set HTTP status to "OK".

	err = json.NewEncoder(w).Encode(hook) // Encode the webhook as JSON and send it.
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleNotifDeleteRequest processes DELETE requests to remove a specific notification webhook by its ID.
func handleNotifDeleteRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")     // Retrieve the ID from the URL path.
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Retrieve token from the URL query parameters.
	if token == "" {            // Check if a token is provided.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.NotificationsID)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token) // Retrieve UUID associated with the API token.
	if UUID == "" {                 // Check if UUID is valid.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.NotificationsID)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" || ID == " " { // Check if the ID is valid.
		log.Printf(constants.ClientConnectNoID, r.RemoteAddr, r.Method, Endpoints.NotificationsID)
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	err := db.DeleteWebhook(ID, UUID) // Delete the specific webhook by ID and UUID from the database.
	if err != nil {
		log.Printf(" %s: Error deleting data from database: %v", r.RemoteAddr, err)
		err := fmt.Sprintf("Error deleting data from database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Set HTTP status to "No Content" upon successful deletion.
}
