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
	ctx    = context.Background()
	Client *auth.Client
)

func init() {
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_FILE"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Panic("Firebase Failed to initialize: ", err)
	}

	Client, err = app.Auth(ctx)
	if err != nil {
		log.Panic("Firebase Failed to initialize Authentication client: ", err)
	}
}
