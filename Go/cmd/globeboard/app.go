package main

import (
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/dashboard"
	"globeboard/internal/handlers/endpoint/util"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"log"
	"net/http"
	"os"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	if !fileExists(constants.FirebaseCredentialPath) {
		log.Fatal("Firebase Credentials file is not mounted")
		return
	}

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
	mux.HandleFunc(Endpoints.RegistrationsID, dashboard.RegistrationsIdHandler)
	mux.HandleFunc(Endpoints.Registrations, dashboard.RegistrationsHandler)
	mux.HandleFunc(Endpoints.Dashboards, dashboard.DashboardsHandler)
	mux.HandleFunc(Endpoints.NotificationsID, dashboard.NotificationsHandler)
	mux.HandleFunc(Endpoints.Notifications, dashboard.NotificationsHandler)
	mux.HandleFunc(Endpoints.Status, dashboard.StatusHandler)

	// Start the HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
