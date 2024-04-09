package util

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"log"
	"net/http"
	"regexp"
)

const (
	ISE = "Internal Server Error"
)

// UserRegistrationHandler handles HTTP GET requests to retrieve supported languages.
func UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		registerUser(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost, http.StatusNotImplemented)
		return
	}
}

func registerUser(w http.ResponseWriter, r *http.Request) {

	// Initialize Firebase
	client, err := authenticate.GetFireBaseAuthClient() // Assuming you have your initFirebase function from earlier
	if err != nil {
		http.Error(w, "Error initializing Firebase Auth", http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate email format
	if !isValidEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if !isValidPassword(password) {
		http.Error(w, "Password does not meet complexity requirements", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	params := (&auth.UserToCreate{}).
		DisplayName(name).
		Email(email).
		Password(password)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")

	UDID := _func.GenerateUID(constants.DocIdLength)
	key := _func.GenerateAPIKey(constants.ApiKeyLength)

	err = db.AddApiKey(UDID, u.UID, key)
	if err != nil {
		log.Printf("error saving API Key: %v\n", err)
		http.Error(w, ISE, http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "Successfully registered user: %s\n"+
		"API Key: %s \n\n"+
		"To delete and get a new API Key utilize: Authorization: \"%v\" at:\n<a href=\"/util/v1/key\">/util/v1/key</a>", u.DisplayName, key, u.UID)
	if err != nil {
		log.Print("Error writing response to ResponseWriter: ", err)
		http.Error(w, ISE, http.StatusInternalServerError)
		return
	}
}

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regex.MatchString(email)
}

func isValidPassword(password string) bool {
	// Check the length
	if len(password) < 12 {
		return false
	}

	// Check for at least one uppercase letter
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return false
	}

	// Check for at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return false
	}

	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$&*]`).MatchString(password)
	if !hasSpecial {
		return false
	}

	// If all checks pass
	return true
}
