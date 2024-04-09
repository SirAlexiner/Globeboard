package util

import (
	"context"
	"encoding/json"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"log"
	"net/http"
)

func APIKeyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleApiKeyGetRequest(w, r)
	case http.MethodDelete:
		handleApiKeyDeleteRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

func handleApiKeyDeleteRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")

	if token == "" {
		http.Error(w, "Please specify API Key to delete: '?token={API_Key}' ", http.StatusBadRequest)
		return
	}

	err := db.DeleteApiKey(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleApiKeyGetRequest(w http.ResponseWriter, r *http.Request) {
	UDID := _func.GenerateUID(constants.DocIdLength)
	key := _func.GenerateAPIKey(constants.ApiKeyLength)

	UUID := r.Header.Get("Authorization")

	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()

	client, err := authenticate.GetFireBaseAuthClient()
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Verify the ID token
	_, err = client.GetUser(ctx, UUID)
	if err != nil {
		log.Printf("error verifying UUID: %v\n", err)
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	err = db.AddApiKey(UDID, UUID, key)
	if err != nil {
		err := fmt.Sprintf("Error creating API Key: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	// Encode books as JSON and send the response.
	if err := json.NewEncoder(w).Encode(key); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
