package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Firestore"
	"globeboard/internal/utils/structs"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	FirebaseClosingErr = "Error closing access to firestore: %v\n"
	IterationFailed    = "failed to iterate over query results: %v\n"
	ParsingError       = "Error parsing document: %v\n"
)

func getFirestoreClient() (*firestore.Client, error) {
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
	client, err := app.Firestore(ctx)
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

func AddApiKey(docID, UUID string, key string) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Printf(FirebaseClosingErr, err)
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.ApiKeyCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	query := ref.Where("UUID", "==", UUID).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err)
		}
		_ = doc
		err = fmt.Errorf("API key is already registered to user")
		return err
	}

	apiKeys := structs.APIKey{
		UUID:   UUID,
		APIKey: key,
	}

	// Create a new document and add it to the firebase
	_, err = ref.Doc(docID).Set(ctx, apiKeys)
	if err != nil {
		err := fmt.Errorf("error saving API key to Database: %v", err)
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
			log.Printf(FirebaseClosingErr, err)
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
			return fmt.Errorf(IterationFailed, err)
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

func GetAPIKeyUUID(apiKey string) string {
	client, err := getFirestoreClient()
	if err != nil {
		log.Printf("Error retriving firestore client: %v", err)
		return ""
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Printf(FirebaseClosingErr, err)
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

	var key structs.APIKey

	// Iterate over the query results
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf(IterationFailed, err)
			return ""
		}
		if err := doc.DataTo(&key); err != nil {
			log.Println("Error parsing document:", err)
			return ""
		}
	}
	clientAuth, err := authenticate.GetFireBaseAuthClient() // Assuming you have your initFirebase function from earlier
	if err != nil {
		log.Printf("Error retriving firebase auth client: %v", err)
		return ""
	}

	_, err = clientAuth.GetUser(ctx, key.UUID)
	if err != nil {
		log.Println("Error getting user:", err)
		return ""
	} else {
		return key.UUID
	}
}

func AddRegistration(docID string, data *structs.CountryInfoGet) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Printf(FirebaseClosingErr, err)
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	// Create a new document and add to it
	_, err = ref.Doc(docID).Set(ctx, data)
	if err != nil {
		log.Println("Error saving data to database" + err.Error())
		return err
	}

	return nil
}

func GetRegistrations(uuid string) ([]*structs.CountryInfoGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Use a context for Firestore operations
	ctx := context.Background()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	// Query all documents
	query := ref.Where("uuid", "==", uuid)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var cis []*structs.CountryInfoGet

	for {
		doc, err := iter.Next()
		var ci *structs.CountryInfoGet
		if errors.Is(err, iterator.Done) {
			break
		}
		if err := doc.DataTo(&ci); err != nil {
			log.Printf(ParsingError, err)
			return nil, err
		}
		cis = append(cis, ci)
	}

	return cis, nil
}

func GetSpecificRegistration(id, uuid string) (*structs.CountryInfoGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Use a context for Firestore operations
	ctx := context.Background()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	query := ref.Where("id", "==", id).Where("uuid", "==", uuid).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var ci *structs.CountryInfoGet

	// Iterate over the query results
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf(IterationFailed, err)
		}
		if err := doc.DataTo(&ci); err != nil {
			log.Println("Error retrieving document:", err)
			// Optionally, continue to the next document instead of returning an error
			// continue
			return nil, err
		}
		return ci, nil
	}

	return nil, errors.New("no document with that id was found")
}

func UpdateRegistration(id, uuid string, data *structs.CountryInfoGet) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Use a context for Firestore operations
	ctx := context.Background()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	query := ref.Where("id", "==", id).Where("uuid", "==", uuid).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	// Iterate over the query results
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err)
		}
		_, err = ref.Doc(doc.Ref.ID).Set(ctx, data)
		if err != nil {
			log.Printf("Error saving data to database: %v\n", err)
			return err
		}
		return nil
	}

	return errors.New("no document with that id was found")
}

func DeleteRegistration(id, uuid string) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Printf(FirebaseClosingErr, err)
			return
		}
	}(client)

	// Create a reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	// Use a context for Firestore operations
	ctx := context.Background()

	query := ref.Where("ID", "==", id).Where("uuid", "==", uuid).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var docID string
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err)
		}
		docID = doc.Ref.ID
	}

	// If docID is empty, the API key does not exist in Firestore
	if docID == "" {
		return fmt.Errorf("ID match was not found")
	}

	// Delete the document with the provided API key
	_, err = ref.Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	fmt.Printf("Registration document %s deleted successfully\n", docID)
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
			log.Printf(FirebaseClosingErr, err)
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
		log.Printf(FirebaseClosingErr, err)
		return err
	}

	return nil
}*/

func GetAllWebhooks() ([]structs.WebhookGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Use a context for Firestore operations
	ctx := context.Background()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	// Query all documents
	docs, err := ref.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error fetching documents: %v\n", err)
		return nil, err
	}

	var webhooks []structs.WebhookGet

	for _, doc := range docs {
		var webhook structs.WebhookGet
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf(ParsingError, err)
			// Optionally, continue to the next document instead of returning an error
			// continue
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

/*func GetWebhooksUser(uuid string) ([]structs.WebhookGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Use a context for Firestore operations
	ctx := context.Background()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	query := ref.Where("uuid", "==", uuid).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var webhooks []structs.WebhookGet

	// Iterate over the query results
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf(IterationFailed, err)
		}
		var webhook structs.WebhookGet
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf(ParsingError, err)
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}*/
