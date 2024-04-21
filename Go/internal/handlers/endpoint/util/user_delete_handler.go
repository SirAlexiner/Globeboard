package util

import (
	"context"
	authenticate "globeboard/auth"
	"log"
	"net/http"
)

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
	if ID == "" || ID == " " {
		http.Error(w, "Please Provide User ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	err := authenticate.Client.DeleteUser(ctx, ID)
	if err != nil {
		log.Printf("error deleting user: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Print("Successfully Deleted User: ", ID)
}
