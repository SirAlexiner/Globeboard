// Package main is the entry point for the application.
package main

import (
	"globeboard/db"
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/dashboard"
	"globeboard/internal/handlers/endpoint/util"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"log"
	"net/http"
	"os"
)

// fileExists checks if a file exists, and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	// Confirm that the Firebase Credentials file is accessible, if not, panic.
	if !fileExists(os.Getenv("FIREBASE_CREDENTIALS_FILE")) {
		log.Panic("Firebase Credentials file is not mounted")
	}
	defer func() {
		// Close the Firestore client connection on application exit
		if err := db.Client.Close(); err != nil {
			log.Printf("Error closing Firestore client: %v", err)
		}
	}()

	// Get the port from the environment variable or set default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Define HTTP endpoints
	mux := http.NewServeMux()
	mux.HandleFunc(Paths.Root, handlers.EmptyHandler)                           // Root endpoint
	mux.HandleFunc(Endpoints.UserRegistration, util.UserRegistrationHandler)    // User registration endpoint
	mux.HandleFunc(Endpoints.UserDeletionID, util.UserDeletionHandler)          // User deletion endpoint
	mux.HandleFunc(Endpoints.ApiKey, util.APIKeyHandler)                        // API key endpoint
	mux.HandleFunc(Endpoints.RegistrationsID, dashboard.RegistrationsIdHandler) // Registrations by ID endpoint
	mux.HandleFunc(Endpoints.Registrations, dashboard.RegistrationsHandler)     // Registrations endpoint
	mux.HandleFunc(Endpoints.DashboardsID, dashboard.DashboardsIdHandler)       // Dashboards by ID endpoint
	mux.HandleFunc(Endpoints.NotificationsID, dashboard.NotificationsIdHandler) // Notifications by ID endpoint
	mux.HandleFunc(Endpoints.Notifications, dashboard.NotificationsHandler)     // Notifications endpoint
	mux.HandleFunc(Endpoints.Status, dashboard.StatusHandler)                   // Status endpoint

	// Start the HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
