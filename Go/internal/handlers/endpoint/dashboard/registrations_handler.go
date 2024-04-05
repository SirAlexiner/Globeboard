// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"net/http"
)

// RegistrationsHandler handles HTTP GET requests to retrieve supported languages.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRegGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is: "+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

// handleGetRequest handles GET requests to retrieve supported languages.
func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	//TODO::Complete HTTP Method Requests
}
