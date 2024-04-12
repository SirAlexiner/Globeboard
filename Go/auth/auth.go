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
	// Use a context for Firebase operations
	ctx = context.Background()
)

func GetFireBaseAuthClient() (*auth.Client, error) {
	// Using the credential file
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_FILE"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println("Credentials not found: " + os.Getenv("FIREBASE_CREDENTIALS_FILE"))
		log.Println("Error on getting the application")
		return nil, err
	}

	//No initial error, so a client is used to gather other information
	client, err := app.Auth(ctx)
	if err != nil {
		// Logging the error
		log.Println("Credentials file: '" + os.Getenv("FIREBASE_CREDENTIALS_FILE") + "' lead to an error.")
		return nil, err
	}

	// No errors, so we return the test client and no error
	return client, nil
}
