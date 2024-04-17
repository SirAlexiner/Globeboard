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
	case http.MethodDelete:
		handleNotifDeleteRequest(w, r)
	case http.MethodPost:
		handleNotifPostRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported.", http.StatusNotImplemented)
		return
	}
}

// handleGetRequest handles GET requests to retrieve readership dashboard for a specific language.
func handleNotifGetRequest(w http.ResponseWriter, r *http.Request) {
	// Check if the link wants to get one webhook or all webhooks, and respond accordingly

	// Make system take JSON

	// Then send ID back
}

func handleNotifDeleteRequest(w http.ResponseWriter, r *http.Request) {
	// Given the id from the url, delete the request

	// Fetch the id from URL

	// Check if webhook with that id exists

	// exists: delete
	// (optional) Ask for deletion

	// Perform deletion
	// Reorder?
}

func handleNotifPostRequest(w http.ResponseWriter, r *http.Request) {

	// Get url

	// Check if url has an additional parameter

	// Has parameter
	//

	// Has no parameter
	// Show all webhooks
}
