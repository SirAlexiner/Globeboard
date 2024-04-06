package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Firestore"
	"globeboard/internal/utils/structs"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"time"
)

const (
	FirebaseClosingErr = "Error closing access to firestore: "
)

func getFirestoreClient(path ...string) (*firestore.Client, error) {
	// Use a service account
	ctx := context.Background()

	// Set the credential path based on if there was given arguments
	var credentialsPath string
	if path != nil {
		credentialsPath = path[0]
	} else {
		credentialsPath = constants.FirebaseCredentialsFilePath
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
	client, err := app.Firestore(ctx)
	if err != nil {
		// Logging the error
		log.Println("Credentials file: '" + credentialsPath + "' lead to an error.")
		return nil, err
	}

	// No errors, so we return the test client and no error
	return client, nil
}

func TestDBConnection() string {
	client, err := getFirestoreClient()
	if err != nil {
		return fmt.Sprintf("%d %s: %v", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Firestore client: %v", err)
		}
	}()

	ctx := context.Background()

	collectionID := "Connectivity"
	documentID := "DB_Connection_Test"

	// Attempt to update a specific document to test connectivity and permissions.
	_, err = client.Collection(collectionID).Doc(documentID).Set(ctx, map[string]interface{}{
		"PSA":         "DO NOT DELETE THIS DOCUMENT!",
		"lastChecked": time.Now(),
	}, firestore.MergeAll)

	if err != nil {
		grpcStatusCode := status.Code(err)
		switch grpcStatusCode {
		case codes.Canceled:
			return fmt.Sprintf("%d %s", http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))
		case codes.DeadlineExceeded:
			return fmt.Sprintf("%d %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
		case codes.PermissionDenied:
			return fmt.Sprintf("%d %s", http.StatusForbidden, http.StatusText(http.StatusForbidden))
		case codes.NotFound:
			// This might indicate the collection or document does not exist,
			//which for this purpose is treated as a connection success
			// since the error was Firestore-specific and not network or permission related.
			return fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK))
		case codes.ResourceExhausted:
			return fmt.Sprintf("%d %s", http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
		case codes.Unauthenticated:
			return fmt.Sprintf("%d %s", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		case codes.Unavailable:
			return fmt.Sprintf("%d %s", http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
		case codes.Unknown, codes.Internal:
			return fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		default:
			// For any other codes, return a generic HTTP 500 error
			return fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}

	// If no error, the document update was successful, indicating good connectivity and permissions.
	return fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK))
}

func AddApiKey(docID string, key string) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Println(FirebaseClosingErr + err.Error())
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.ApiKeyCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	apiKeys := structs.APIKey{APIKey: key}

	// Create a new document and add it to the firebase
	_, err = ref.Doc(docID).Set(ctx, apiKeys)
	if err != nil {
		errString := fmt.Sprintf("Error saving API key to Database: %v", err)
		err := errors.New(errString)
		return err
	}

	return nil
}

func DeleteApiKey(apiKey string) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Println(FirebaseClosingErr + err.Error())
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.ApiKeyCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	query := ref.Where("APIKey", "==", apiKey).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var docID string
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate over query results: %v", err)
		}
		docID = doc.Ref.ID
	}

	// If docID is empty, the API key does not exist in Firestore
	if docID == "" {
		return fmt.Errorf("API key not found")
	}

	// Delete the document with the provided API key
	_, err = ref.Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	fmt.Printf("API key %s deleted successfully\n", apiKey)
	return nil
}

func DoesAPIKeyExists(apiKey string) (bool, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return false, err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Println(FirebaseClosingErr + err.Error())
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.ApiKeyCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	// Query Firestore to check if the API key exists
	query := ref.Where("APIKey", "==", apiKey).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	// Iterate over the query results
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return false, fmt.Errorf("failed to iterate over query results: %v", err)
		}
		// If a document is found, the API key exists in Firestore
		_ = doc // You can process the document if needed
		return true, nil
	}

	// If no matching document is found, the API key does not exist in Firestore
	return false, nil
}

func AddRegistration(docID, userID string, data *structs.CountryInfoPost) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Println(FirebaseClosingErr + err.Error())
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	// Create a new document and add to it
	_, err = ref.Doc(docID).Set(ctx, map[string]interface{}{
		"id":         userID,
		"country":    data.Country,
		"isoCode":    data.IsoCode,
		"features":   data.Features,
		"lastChange": time.Now(),
	})
	if err != nil {
		log.Println(FirebaseClosingErr + err.Error())
		return err
	}

	return nil
}

/*
func AddWebhook(userID, docID string, webhook structs.WebhookPost) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Println(FirebaseClosingErr + err.Error())
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	// Create a new document and add it to the
	_, err = ref.Doc(docID).Set(ctx, map[string]interface{}{
		"id":         userID,
	    "url": 		  webhook.URL,
   		"country":    webhook.Country,
   		"event":      webhook.Event,
		"lastChange": time.Now(),
	})
	if err != nil {
		log.Println(FirebaseClosingErr + err.Error())
		return err
	}

	return nil
}*/

func GetWebhooks() ([]structs.WebhookGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Println(FirebaseClosingErr + err.Error())
		}
	}()

	// Use a context for Firestore operations
	ctx := context.Background()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	// Query all documents
	docs, err := ref.Documents(ctx).GetAll()
	if err != nil {
		log.Println("Error fetching documents:", err)
		return nil, err
	}

	var webhooks []structs.WebhookGet

	for _, doc := range docs {
		var webhook structs.WebhookGet
		if err := doc.DataTo(&webhook); err != nil {
			log.Println("Error parsing document:", err)
			// Optionally, continue to the next document instead of returning an error
			// continue
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}
