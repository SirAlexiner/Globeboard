// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleNotifPostRequest(w, r)
	case http.MethodGet:
		handleNotifGetAllRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet+"\n"+http.MethodPatch, http.StatusNotImplemented)
		return
	}
}

func handleNotifPostRequest(w http.ResponseWriter, r *http.Request) {
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

	if r.Body == nil {
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	var webhook *structs.WebhookGet
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		err := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	UDID := _func.GenerateUID(constants.DocIdLength)
	ID := _func.GenerateUID(constants.IdLength)

	webhook.ID = ID
	webhook.UUID = UUID

	err := db.AddWebhook(UDID, webhook)
	if err != nil {
		log.Println("Error saving data to database" + err.Error())
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	hook, err := db.GetSpecificWebhook(ID, UUID)
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error confirming data added to database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": hook.ID,
	}

	w.Header().Set(ContentType, ApplicationJSON)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleNotifGetAllRequest(w http.ResponseWriter, r *http.Request) {
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
	regs, err := db.GetWebhooksUser(UUID)
	if err != nil {
		errmsg := fmt.Sprint("Error retrieving webhooks from database: ", err)
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON)

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(regs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
