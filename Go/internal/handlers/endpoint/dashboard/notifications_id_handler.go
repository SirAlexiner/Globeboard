// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	"log"
	"net/http"
)

// NotificationsIdHandler handles requests to retrieve readership dashboard for a specific language.
func NotificationsIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNotifGetRequest(w, r)
	case http.MethodDelete:
		handleNotifDeleteRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet+"\n"+http.MethodPatch, http.StatusNotImplemented)
		return
	}
}

func handleNotifGetRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token)
	if UUID == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	hook, err := db.GetSpecificWebhook(ID, UUID)
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error retrieving data from database", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set(ContentType, ApplicationJSON)

	// Write the status code to the response
	w.WriteHeader(http.StatusOK)

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(hook)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleNotifDeleteRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token)
	if UUID == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	err := db.DeleteWebhook(ID, UUID)
	if err != nil {
		err := fmt.Sprintf("Error deleting data from database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
