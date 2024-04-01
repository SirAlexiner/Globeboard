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
	"strconv"
	"time"
)

// getEndpointStatus returns the HTTP status code of the provided endpoint.
func getEndpointStatus(endpointURL string) string {
	// Create new request
	r, err := http.NewRequest(http.MethodGet, endpointURL, nil)
	if err != nil {
		// Log and handle the error if request creation fails.
		err := fmt.Errorf("error in creating request: %s", err.Error())
		log.Println(err)
	}

	// Set content type header
	r.Header.Add("content-type", "application/json")

	// Create HTTP client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		// Log and handle the error if request execution fails.
		err := fmt.Errorf("error in response: %s", err.Error())
		log.Println(err)
	}

	return res.Status
}

// StatusHandler handles requests to retrieve the status of Endpoints.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r)
	default:
		// If method is not supported, return an error response.
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

// handleStatusGetRequest handles GET requests to retrieve the status of Endpoints.
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Please provide API Token", http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf("Error checking API key: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	if !exists {
		err := fmt.Sprintf("API key not accepted")
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	// Initialize a status response.
	status := structs.StatusResponse{
		CountriesApi:   getEndpointStatus(External.CountriesAPI + "all"),
		MeteoApi:       getEndpointStatus(External.OpenMeteoAPI),
		CurrencyApi:    getEndpointStatus(External.CurrencyAPI + "nok"),
		FirebaseDB:     db.TestDBConnection(),
		NotificationDb: strconv.Itoa(http.StatusNotImplemented) + " Not Implemented", // TODO::Update with Notification DB
		Webhooks:       0,                                                            //TODO::Get Actual number of webhooks
		Version:        constants.APIVersion,
		// Calculate uptime since the last restart of the service.
		UptimeInSeconds: fmt.Sprintf("%f Seconds", time.Since(startTime).Seconds()),
	}

	// Set content type header
	w.Header().Add("content-type", "application/json")

	// Encode status as JSON and send the response.
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		// If encoding fails, return an error response.
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// startTime keeps track of the service start time.
var startTime = time.Now()
