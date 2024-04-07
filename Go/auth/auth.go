package authenticate

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"globeboard/internal/utils/constants"
	"google.golang.org/api/option"
	"log"
	"os"
)

func GetFireBaseAuthClient() (*auth.Client, error) {
	// Use a service account
	ctx := context.Background()

	// Set the credential path based on if it is running in Docker
	var credentialsPath string
	if runningInDocker() {
		credentialsPath = constants.FirebaseCredentialsDockerPath
	} else {
		credentialsPath = constants.FirebaseCredentialsDefaultPath
	}

	// Using the credential file
	sa := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println("Credentials not found: " + credentialsPath)
		log.Println("Error on getting the application")
		return nil, err
	}

	//No initial error, so a client is used to gather other information
	client, err := app.Auth(ctx)
	if err != nil {
		// Logging the error
		log.Println("Credentials file: '" + credentialsPath + "' lead to an error.")
		return nil, err
	}

	// No errors, so we return the test client and no error
	return client, nil
}

// runningInDocker checks if the application is running inside a Docker container.
func runningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true // .dockerenv exists
	}
	return false // .dockerenv does not exist
}
