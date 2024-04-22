// Package dashboard provides handlers for managing dashboard-related functionalities through HTTP endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
)

// NotificationsHandler handles HTTP requests related to notification webhooks.
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost: // Handle POST request
		handleNotifPostRequest(w, r)
	case http.MethodGet: // Handle GET request
		handleNotifGetAllRequest(w, r)
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.Notifications, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

// handleNotifPostRequest processes POST requests to create a new notification webhook.
func handleNotifPostRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Retrieve token from the URL query parameters.
	if token == "" {            // Check if a token is provided.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.Notifications)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(r.RemoteAddr, token) // Retrieve UUID associated with the API token.
	if UUID == "" {                               // Check if UUID is valid.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.Notifications)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}

	if r.Body == nil { // Check if request body is empty.
		log.Printf(constants.ClientConnectEmptyBody, r.RemoteAddr, r.Method, Endpoints.Notifications)
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	var webhook *structs.WebhookInternal
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil { // Decode the JSON request body into webhook struct.
		err := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	UDID := _func.GenerateUID(constants.DocIdLength) // Generate a unique document ID.
	ID := _func.GenerateUID(constants.IdLength)      // Generate a unique ID for the webhook.

	webhook.ID = ID
	webhook.UUID = UUID

	err := db.AddWebhook(r.RemoteAddr, UDID, webhook) // Add the webhook to the database.
	if err != nil {
		log.Println("Error saving data to database" + err.Error())
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	hook, err := db.GetSpecificWebhook(r.RemoteAddr, ID, UUID) // Retrieve the newly added webhook to confirm its addition.
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error confirming data added to database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": hook.ID, // Prepare response data with the new webhook ID.
	}

	w.Header().Set(ContentType, ApplicationJSON) // Set the content type of the response.

	w.WriteHeader(http.StatusCreated) // Set HTTP status to "Created".

	err = json.NewEncoder(w).Encode(response) // Encode the response as JSON and send it.
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleNotifGetAllRequest processes GET requests to retrieve all notification webhooks for a user.
func handleNotifGetAllRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Retrieve token from the URL query parameters.
	if token == "" {            // Check if a token is provided.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.Notifications)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(r.RemoteAddr, token) // Retrieve UUID associated with the API token.
	if UUID == "" {                               // Check if UUID is valid.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.Notifications)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	regs, err := db.GetWebhooksUser(r.RemoteAddr, UUID) // Retrieve all webhooks associated with the user (UUID).
	if err != nil {
		log.Printf("%s: Error retrieving webhooks from database: %v", r.RemoteAddr, err)
		errmsg := fmt.Sprint("Error retrieving webhooks from database: ", err)
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON) // Set the content type of the response.

	w.WriteHeader(http.StatusOK) // Set HTTP status to "OK".

	err = json.NewEncoder(w).Encode(regs) // Encode the webhooks as JSON and send it.
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
