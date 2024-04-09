package authenticate

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"globeboard/internal/utils/constants"
	"google.golang.org/api/option"
	"log"
)

var (
	// Use a context for Firebase operations
	ctx = context.Background()
)

func GetFireBaseAuthClient() (*auth.Client, error) {
	// Using the credential file
	sa := option.WithCredentialsFile(constants.FirebaseCredentialPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println("Credentials not found: " + constants.FirebaseCredentialPath)
		log.Println("Error on getting the application")
		return nil, err
	}

	//No initial error, so a client is used to gather other information
	client, err := app.Auth(ctx)
	if err != nil {
		// Logging the error
		log.Println("Credentials file: '" + constants.FirebaseCredentialPath + "' lead to an error.")
		return nil, err
	}

	// No errors, so we return the test client and no error
	return client, nil
}
