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

	UUID := r.Header.Get("Authorization")

	ctx := context.Background()

	_, err := authenticate.Client.GetUser(ctx, UUID)
	if err != nil {
		log.Printf("error verifying UUID: %v\n", err)
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	if token == "" {
		http.Error(w, "Please specify API Key to delete: '?token={API_Key}' ", http.StatusBadRequest)
		return
	}

	err = db.DeleteApiKey(UUID, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleApiKeyGetRequest(w http.ResponseWriter, r *http.Request) {
	UDID := _func.GenerateUID(constants.DocIdLength)
	key := _func.GenerateAPIKey(constants.ApiKeyLength)

	UUID := r.Header.Get("Authorization")

	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()

	_, err := authenticate.Client.GetUser(ctx, UUID)
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

	response := map[string]string{
		"token": key,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
