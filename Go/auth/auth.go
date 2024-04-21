// Package authenticate provides functionality for initializing and accessing Firebase Authentication.
package authenticate

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	ctx    = context.Background() // Background context for Firebase operations
	Client *auth.Client           // Singleton Firebase Authentication client
)

func init() {
	// Load Firebase service account credentials from environment variable
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_FILE"))

	// Initialize Firebase app with the loaded credentials
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Panic("Firebase Failed to initialize: ", err)
	}

	// Initialize the Firebase Authentication client
	Client, err = app.Auth(ctx)
	if err != nil {
		log.Panic("Firebase Failed to initialize Authentication client: ", err)
	}
}
