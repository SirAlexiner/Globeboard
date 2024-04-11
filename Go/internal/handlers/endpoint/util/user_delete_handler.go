package util

import (
	"context"
	authenticate "globeboard/auth"
	"log"
	"net/http"
)

// UserDeletionHandler handles HTTP Delete requests
func UserDeletionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost, http.StatusNotImplemented)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
	if ID == "" {
		http.Error(w, "Please Provide User ID", http.StatusBadRequest)
		return
	}

	// Initialize Firebase
	client, err := authenticate.GetFireBaseAuthClient() // Assuming you have your initFirebase function from earlier
	if err != nil {
		http.Error(w, "Error initializing Firebase Auth", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	err = client.DeleteUser(ctx, ID)
	if err != nil {
		log.Printf("error deleting user: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
