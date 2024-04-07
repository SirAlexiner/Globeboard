package main

import (
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/dashboard"
	"globeboard/internal/handlers/endpoint/util"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get the port from the environment variable or set default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Define HTTP endpoints
	mux := http.NewServeMux()
	mux.HandleFunc(Paths.Root, handlers.EmptyHandler)
	mux.HandleFunc(Endpoints.UserRegistration, util.UserRegistrationHandler)
	mux.HandleFunc(Endpoints.ApiKey, util.APIKeyHandler)
	mux.HandleFunc(Endpoints.RegistrationsID, dashboard.RegistrationsHandler)
	mux.HandleFunc(Endpoints.Registrations, dashboard.RegistrationsHandler)
	mux.HandleFunc(Endpoints.Dashboards, dashboard.DashboardsHandler)
	mux.HandleFunc(Endpoints.NotificationsID, dashboard.NotificationsHandler)
	mux.HandleFunc(Endpoints.Notifications, dashboard.NotificationsHandler)
	mux.HandleFunc(Endpoints.Status, dashboard.StatusHandler)

	// Start the HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
