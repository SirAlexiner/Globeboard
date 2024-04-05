// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"net/http"
)

// NotificationsHandler handles requests to retrieve readership dashboard for a specific language.
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNotifGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
		return
	}
}

// handleGetRequest handles GET requests to retrieve readership dashboard for a specific language.
func handleNotifGetRequest(w http.ResponseWriter, r *http.Request) {
	//TODO::Complete HTTP Method Requests
}
