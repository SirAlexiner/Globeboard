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

func AddNewWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Add some webhook and return some id
	} else {
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
	}
}

func DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		// Delete a registered webhook
	} else {
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
	}
}

func ViewOneWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Show one webhook
	} else {
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
	}
}

func ViewAllWebhooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Show all webhooks
	} else {
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
	}
}
