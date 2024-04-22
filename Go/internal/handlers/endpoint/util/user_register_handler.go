// Package util provides HTTP handlers for user and API key management within the application.
package util

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	authenticate "globeboard/auth"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"log"
	"net/http"
	"regexp"
)

const (
	ISE = "Internal Server Error" // ISE is the error message used when an internal server error occurs.
)

// UserRegistrationHandler handles HTTP requests for user registration.
func UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		registerUser(w, r) // Handle POST requests
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.UserRegistration, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost, http.StatusNotImplemented)
		return
	}
}

// registerUser processes the user registration, including input validation and user creation.
func registerUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")     // Extract username from form data.
	email := r.FormValue("email")       // Extract email from form data.
	password := r.FormValue("password") // Extract password from form data.

	if !isValidEmail(email) { // Validate the email.
		// Log a message indicating that a client attempted to register using a malformed email.
		log.Printf("%s attempted to register a user with malformed email.", r.RemoteAddr)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !isValidPassword(password) { // Validate password strength.
		// Log a message indicating that a client attempted to register using a weak password.
		log.Printf("%s attempted to register a user with a weak password.", r.RemoteAddr)
		http.Error(w, "Password does not meet complexity requirements", http.StatusBadRequest)
		return
	}

	ctx := context.Background() // Create a new background context.

	// Define user creation parameters.
	params := (&auth.UserToCreate{}).
		DisplayName(name).
		Email(email).
		Password(password)

	u, err := authenticate.Client.CreateUser(ctx, params) // Attempt to create user in Firebase.
	if err != nil {
		log.Printf("%s: Error creating user: %v\n", r.RemoteAddr, err) // Log the error.
		http.Error(w, err.Error(), http.StatusInternalServerError)     // Report creation error.
		return
	}

	w.Header().Set("content-type", "application/json") // Set response content type.

	UDID := _func.GenerateUID(constants.DocIdLength)    // Generate a unique document ID.
	key := _func.GenerateAPIKey(constants.ApiKeyLength) // Generate a new API key.

	err = db.AddApiKey(UDID, u.UID, key) // Store the new API key in the database.
	if err != nil {
		log.Printf("%s Error saving API Key: %v\n", r.RemoteAddr, err) // Log the error.
		http.Error(w, ISE, http.StatusInternalServerError)             // Report API key storage error.
		return
	}

	// Prepare the JSON response with the user ID and token.
	response := struct {
		Token  string `json:"token"`  // API token.
		UserID string `json:"userid"` // Firebase user ID.
	}{
		Token:  key,
		UserID: u.UID,
	}

	w.WriteHeader(http.StatusCreated) // Set HTTP status to 201 Created.

	encoder := json.NewEncoder(w)  // Initialize JSON encoder.
	err = encoder.Encode(response) // Encode the response as JSON.
	if err != nil {
		log.Printf("%s: Error encoding JSON response: %v\n", r.RemoteAddr, err) // Log encoding error.
		http.Error(w, ISE, http.StatusInternalServerError)                      // Report encoding error.
		return
	}

	// Log successful creation events.
	log.Printf("%s: Successfully created user: %v with API Key: %v\n", r.RemoteAddr, response.UserID, response.Token)
}

// isValidEmail checks if the provided email string matches the expected format.
func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // Regular expression for email validation.
	return regex.MatchString(email)                                               // Return true
	// if email matches the regex.
}

// isValidPassword checks if the provided password meets complexity requirements.
func isValidPassword(password string) bool {
	if len(password) < 12 { // Password length should be at least 12 characters.
		return false
	}
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)  // Check for at least one uppercase letter.
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)      // Check for at least one digit.
	hasSpecial := regexp.MustCompile(`[!@#$&*]`).MatchString(password) // Check for at least one special character.
	return hasUppercase && hasDigit && hasSpecial                      // Return true if all conditions are met.
}
