// Package dashboard provides handlers for managing dashboard-related functionalities through HTTP endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
	"time"
)

// getEndpointStatus sends a GET request to the specified endpoint URL and returns the HTTP status code as a string.
func getEndpointStatus(endpointURL string) string {
	r, err := http.NewRequest(http.MethodGet, endpointURL, nil) // Create a new GET request for the endpoint URL.
	if err != nil {
		log.Printf("error in creating request: %v", err)
		return "Failed to create request"
	}

	r.Header.Add("content-type", "application/json") // Set content-type of request.

	client := &http.Client{Timeout: 10 * time.Second} // Initialize a new HTTP client with a timeout.
	defer client.CloseIdleConnections()               // Ensure that idle connections are closed upon function exit.

	res, err := client.Do(r) // Execute the request.
	if err != nil {
		log.Printf("error in receiving response: %v", err)
		return "Failed to connect"
	}
	// Ensure the response body is closed after the function exits, checking for errors.
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	return res.Status // Return the status code of the response.
}

// StatusHandler routes requests based on HTTP method to handle status retrieval.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r) // Handle GET requests with handleStatusGetRequest.
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, Endpoints.Status, r.Method)
		http.Error(w, fmt.Sprintf("REST Method: %s not supported. Only GET is supported for this endpoint", r.Method), http.StatusNotImplemented)
	}
}

// handleStatusGetRequest processes GET requests to retrieve and report the status of various services and endpoints.
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token") // Retrieve the API token from query parameters.
	if token == "" {                    // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.Method, Endpoints.Status)
		http.Error(w, "Please provide API Token", http.StatusUnauthorized)
		return
	}

	UUID := db.GetAPIKeyUUID(token) // Retrieve the UUID associated with the API token.
	if UUID == "" {                 // Validate UUID presence.
		log.Printf(constants.ClientConnectUnauthorized, r.Method, Endpoints.Status)
		http.Error(w, "API key not accepted", http.StatusNotAcceptable)
		return
	}

	webhooksUser, err := db.GetWebhooksUser(UUID) // Retrieve user data associated with webhooks.
	if err != nil {
		log.Printf("Error retrieving user's webhooks: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a status response object.
	status := structs.StatusResponse{
		CountriesApi:    getEndpointStatus(External.CountriesAPI + "alpha?codes=no"),
		MeteoApi:        getEndpointStatus(External.OpenMeteoAPI),
		CurrencyApi:     getEndpointStatus(External.CurrencyAPI + "nok"),
		FirebaseDB:      db.TestDBConnection(), // Test the database connection.
		Webhooks:        len(webhooksUser),
		Version:         constants.APIVersion,                                       // Include the API version.
		UptimeInSeconds: fmt.Sprintf("%f Seconds", time.Since(startTime).Seconds()), // Calculate uptime.
	}

	w.Header().Set("Content-Type", "application/json") // Set response content type to application/json.

	err = json.NewEncoder(w).Encode(status) // Encode the status response to JSON and send it.
	if err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("Error during encoding: %v", err), http.StatusInternalServerError)
		return
	}
}

var startTime = time.Now() // Track the start time of the application.
