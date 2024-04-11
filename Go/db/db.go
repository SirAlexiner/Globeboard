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
)

const (
	FirebaseClosingErr = "Error closing access to firestore: %v\n"
	IterationFailed    = "failed to iterate over query results: %v\n"
	ParsingError       = "Error parsing document: %v\n"
)

var (
	// Use a context for Firestore operations
	ctx = context.Background()
)

func getFirestoreClient() (*firestore.Client, error) {

	// Using the credential file
	sa := option.WithCredentialsFile(constants.FirebaseCredentialPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println("Credentials not found: " + constants.FirebaseCredentialPath)
		log.Println("Error on getting the application")
		return nil, err
	}

	//No initial error, so a client is used to gather other information
	client, err := app.Firestore(ctx)
	if err != nil {
		// Logging the error
		log.Println("Credentials file: '" + constants.FirebaseCredentialPath + "' lead to an error.")
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

	collectionID := "Connectivity"
	documentID := "DB_Connection_Test"

	// Attempt to update a specific document to test connectivity and permissions.
	_, err = client.Collection(collectionID).Doc(documentID).Set(ctx, map[string]interface{}{
		"PSA":         "DO NOT DELETE THIS DOCUMENT!",
		"lastChecked": firestore.ServerTimestamp,
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

	iter := ref.Where("UUID", "==", UUID).Limit(1).Documents(ctx)
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
		err = errors.New("API key is already registered to user")
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

	iter := ref.Where("APIKey", "==", apiKey).Limit(1).Documents(ctx)
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
		return errors.New("API key not found")
	}

	// Delete the document with the provided API key
	_, err = ref.Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete API Key: %v", err)
	}

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

	// Query Firestore to check if the API key exists
	iter := ref.Where("APIKey", "==", apiKey).Limit(1).Documents(ctx)
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

	// Create a new document and add to it
	_, err = ref.Doc(docID).Set(ctx, map[string]interface{}{
		"ID":         data.ID,
		"UUID":       data.UUID,
		"Country":    data.Country,
		"IsoCode":    data.IsoCode,
		"Features":   data.Features,
		"Lastchange": firestore.ServerTimestamp,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetRegistrations(UUID string) ([]*structs.CountryInfoGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	// Query all documents
	docs, _ := ref.Where("UUID", "==", UUID).OrderBy("Lastchange", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error fetching Registration: %v\n", err)
		return nil, err
	}

	var cis []*structs.CountryInfoGet

	for _, doc := range docs {
		var ci *structs.CountryInfoGet
		if err := doc.DataTo(&ci); err != nil {
			log.Printf(ParsingError, err)
			return nil, err
		}
		cis = append(cis, ci)
	}

	return cis, nil
}

func GetSpecificRegistration(ID, UUID string) (*structs.CountryInfoGet, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
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
			return nil, err
		}
		return ci, nil
	}

	return nil, errors.New("no registration with that ID was found")
}

func UpdateRegistration(ID, UUID string, data *structs.CountryInfoGet) error {
	client, err := getFirestoreClient()
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.RegistrationCollection)

	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
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
		_, err = ref.Doc(doc.Ref.ID).Set(ctx, map[string]interface{}{
			"ID":         data.ID,
			"UUID":       data.UUID,
			"Country":    data.Country,
			"IsoCode":    data.IsoCode,
			"Features":   data.Features,
			"Lastchange": firestore.ServerTimestamp,
		})
		if err != nil {
			log.Printf("Error saving data to database: %v\n", err)
			return err
		}
		return nil
	}

	return errors.New("no registration  with that ID was found")
}

func DeleteRegistration(ID, UUID string) error {
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

	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
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
		return errors.New("ID match was not found")
	}

	// Delete the document with the provided API key
	_, err = ref.Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	log.Printf("Registration document %s deleted successfully\n", docID)
	return nil
}

func AddWebhook(docID string, webhook *structs.WebhookGet) error {
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

	// Create a new document and add it to the
	_, err = ref.Doc(docID).Set(ctx, webhook)
	if err != nil {
		log.Printf(FirebaseClosingErr, err)
		return err
	}

	return nil
}

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

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	// Query all documents
	docs, err := ref.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error fetching all stored Webhooks: %v\n", err)
		return nil, err
	}

	var webhooks []structs.WebhookGet

	for _, doc := range docs {
		var webhook structs.WebhookGet
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf(ParsingError, err)
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

func GetWebhooksUser(UUID string) ([]structs.WebhookResponse, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	// Query all documents
	docs, err := ref.Where("UUID", "==", UUID).Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error fetching users webhooks: %v\n", err)
		return nil, err
	}

	var webhooks []structs.WebhookResponse

	for _, doc := range docs {
		var webhook structs.WebhookResponse
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf(ParsingError, err)
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

func GetSpecificWebhook(ID, UUID string) (*structs.WebhookResponse, error) {
	client, err := getFirestoreClient()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf(FirebaseClosingErr, err)
		}
	}()

	// Reference to the Firestore collection
	ref := client.Collection(Firestore.WebhookCollection)

	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop()

	var webhook *structs.WebhookResponse

	// Iterate over the query results
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf(IterationFailed, err)
		}
		if err := doc.DataTo(&webhook); err != nil {
			log.Println("Error retrieving document:", err)
			return nil, err
		}
		return webhook, nil
	}

	return nil, errors.New("no document with that ID was found")
}

func DeleteWebhook(ID, UUID string) error {
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

	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
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

	log.Printf("Webhook %s deleted successfully\n\n", docID)
	return nil
}
