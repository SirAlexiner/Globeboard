package library

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"net/http"
)

func APIKeyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleApiKeyGetRequest(w, r)
	case http.MethodDelete:
		handleApiKeyDeleteRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
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
	UUID := _func.GenerateUUID(constants.UserIdLength)
	key := _func.GenerateAPIKey(constants.ApiKeyLength)

	// Your username and password
	username := "GutendexAdmin"
	password := "4dm1n_4cc355_6r4n73d!"

	// Concatenate username and password with a colon
	auth := username + ":" + password

	// Base64 encode the username and password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	pass := r.Header.Get("Authorization")

	w.Header().Set("Content-Type", "application/json")

	if pass != encodedAuth {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	err := db.AddApiKey(UUID, key)
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
