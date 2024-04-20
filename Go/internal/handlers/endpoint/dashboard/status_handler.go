// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
	"time"
)

func getEndpointStatus(endpointURL string) string {
	r, err := http.NewRequest(http.MethodGet, endpointURL, nil)
	if err != nil {
		err := fmt.Errorf("error in creating request: %v", err)
		log.Println(err)
	}

	r.Header.Add("content-type", "application/json")

	client := &http.Client{}
	defer client.CloseIdleConnections()

	res, err := client.Do(r)
	if err != nil {
		err := fmt.Errorf("error in response: %v", err)
		log.Println(err)
	}

	return res.Status
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Please provide API Token", http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token)
	if UUID == "" {
		err := fmt.Sprintf("API key not accepted")
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}

	webhooksUser, err := db.GetWebhooksUser(UUID)
	if err != nil {
		log.Print("Error retrieving users webhooks:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	status := structs.StatusResponse{
		CountriesApi:    getEndpointStatus(External.CountriesAPI + "alpha?codes=no"),
		MeteoApi:        getEndpointStatus(External.OpenMeteoAPI),
		CurrencyApi:     getEndpointStatus(External.CurrencyAPI + "nok"),
		FirebaseDB:      db.TestDBConnection(),
		Webhooks:        len(webhooksUser),
		Version:         constants.APIVersion,
		UptimeInSeconds: fmt.Sprintf("%f Seconds", time.Since(startTime).Seconds()),
	}

	w.Header().Add("content-type", "application/json")

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

var startTime = time.Now()
