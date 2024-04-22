// Package db provides data access functions for interacting with Firestore.
package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/internal/utils/constants/Firestore"
	"globeboard/internal/utils/structs"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"os"
)

const (
	// IterationFailed Error for when iteration over Firestore results fails.
	IterationFailed = "failed to iterate over query results: %v\n"
)

var (
	ctx    = context.Background() // Global context for Firestore operations, used across all Firestore calls.
	Client *firestore.Client      // Singleton Firestore client.
	err    error                  // Variable to handle errors globally within this package.
)

// init initializes the Firestore client using environment variables for credentials and project ID.
func init() {
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_FILE"))      // Set up the credential file from environment.
	Client, err = firestore.NewClient(ctx, os.Getenv("FIRESTORE_PROJECT_ID"), sa) // Create a new Firestore client.
	if err != nil {
		log.Panic("Firestore was unable to initialize: ", err) // Panic if Firestore client initialization fails.
	}
}

// TestDBConnection tests the Firestore database connection by attempting to write and immediately read a document.
func TestDBConnection() string {
	collectionID := "Connectivity"     // Define the collection ID for connection tests.
	documentID := "DB_Connection_Test" // Define the document ID for connection tests.

	// Attempt to set a document in the Firestore collection, tagging it with server timestamp.
	_, err = Client.Collection(collectionID).Doc(documentID).Set(ctx, map[string]interface{}{
		"PSA":         "DO NOT DELETE THIS DOCUMENT!",
		"lastChecked": firestore.ServerTimestamp,
	}, firestore.MergeAll)

	// Handle potential errors and map them to HTTP status codes.
	grpcStatusCode := status.Code(err)
	switch grpcStatusCode {
	case codes.OK:
		return fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK))
	case codes.Canceled:
		return fmt.Sprintf("%d %s", http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))
	case codes.DeadlineExceeded:
		return fmt.Sprintf("%d %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
	case codes.PermissionDenied:
		return fmt.Sprintf("%d %s", http.StatusForbidden, http.StatusText(http.StatusForbidden))
	case codes.NotFound:
		// Treat not found as OK for this operation; it indicates collection/document was simply not found.
		// Another error in-of-itself.
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
		return fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

// AddApiKey adds a new API key to Firestore, ensuring it does not already exist for the provided user (UUID).
func AddApiKey(docID, UUID string, key string) error {
	ref := Client.Collection(Firestore.ApiKeyCollection) // Reference to the APIKey collection in Firestore.

	// Query for existing API keys with the same UUID.
	iter := ref.Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	for {
		_, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		err = errors.New("API key is already registered to user")
		return err // Return error if an existing key is found.
	}

	apiKeys := structs.APIKey{UUID: UUID, APIKey: key} // Create an APIKey struct to be saved.

	_, err = ref.Doc(docID).Set(ctx, apiKeys) // Set the APIKey document in Firestore.
	if err != nil {
		err := fmt.Errorf("error saving API key to Database: %v", err)
		return err // Return formatted error if setting the document fails.
	}

	log.Printf("API key %s created successfully.", apiKeys.APIKey) // Log success.
	return nil                                                     // Return nil error on success.
}

// DeleteApiKey deletes an API key from Firestore based on UUID and key value.
func DeleteApiKey(UUID, apiKey string) error {
	ref := Client.Collection(Firestore.ApiKeyCollection) // Reference to the APIKey collection in Firestore.

	// Query for the API key document based on UUID and key.
	iter := ref.Where("UUID", "==", UUID).Where("APIKey", "==", apiKey).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	var docID string // Variable to store the document ID of the found API key.
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		docID = doc.Ref.ID // Store the document ID.
	}

	if docID == "" {
		return errors.New("API key not found") // Return error if no document ID was found.
	} // Handle if the API key was not found.

	_, err = ref.Doc(docID).Delete(ctx) // Delete the document from Firestore.
	if err != nil {
		return fmt.Errorf("failed to delete API Key: %v", err) // Return formatted error if delete fails.
	}

	log.Printf("API key %s deleted successfully.", apiKey) // Log success.
	return nil                                             // Return nil error on success.
}

// GetAPIKeyUUID retrieves the UUID associated with a specific API key from Firestore.
func GetAPIKeyUUID(apiKey string) string {
	ref := Client.Collection(Firestore.ApiKeyCollection) // Reference to the APIKey collection in Firestore.

	// Query for the API key document based on the key value.
	iter := ref.Where("APIKey", "==", apiKey).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	var key structs.APIKey // Variable to store the API key data.
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			log.Printf(IterationFailed, err)
			return "" // Return an empty string on error.
		}
		if err := doc.DataTo(&key); err != nil {
			log.Println("Error parsing document:", err)
			return "" // Return an empty string on parsing error.
		}
	}

	_, err = authenticate.Client.GetUser(ctx, key.UUID) // Authenticate the user based on UUID.
	if err != nil {
		log.Println("Error getting user:", err)
		return "" // Return an empty string if user authentication fails.
	} else {
		log.Printf("UUID: %s successfully retrieved from API key: %s.", key.UUID, key.APIKey)
		return key.UUID // Return the UUID on success.
	}
}

// AddRegistration adds a new registration document to Firestore.
func AddRegistration(docID string, data *structs.CountryInfoInternal) error {
	ref := Client.Collection(Firestore.RegistrationCollection) // Reference to the Registration collection.

	// Set the registration document in Firestore with the given ID and data.
	_, err = ref.Doc(docID).Set(ctx, map[string]interface{}{
		"ID":         data.ID,
		"UUID":       data.UUID,
		"Country":    data.Country,
		"IsoCode":    data.IsoCode,
		"Features":   data.Features,
		"Lastchange": firestore.ServerTimestamp, // Use server timestamp to record last change.
	})
	if err != nil {
		return err // Return error if the document set operation fails.
	}

	log.Printf("Registration documents %s created successfully.", data.ID)
	return nil // Return nil if the addition is successful.
}

// GetRegistrations retrieves all registration documents for a given user (UUID) from Firestore.
func GetRegistrations(UUID string) ([]*structs.CountryInfoInternal, error) {
	ref := Client.Collection(Firestore.RegistrationCollection) // Reference to the Registration collection.

	// Query and retrieve all documents where 'UUID' matches, ordered by 'Lastchange' descending.
	docs, err := ref.Where("UUID", "==", UUID).OrderBy("Lastchange", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err // Return error if the fetch operation fails.
	}

	var cis []*structs.CountryInfoInternal // Slice to store the fetched documents.

	for _, doc := range docs {
		var ci *structs.CountryInfoInternal
		if err := doc.DataTo(&ci); err != nil {
			return nil, err // Return error if parsing any document fails.
		}
		cis = append(cis, ci) // Append the parsed document to the slice.
	}
	log.Printf("Registration documents for user: %s retrieved successfully.", UUID)
	return cis, nil // Return the slice of documents.
}

// GetSpecificRegistration retrieves a specific registration document by ID and UUID from Firestore.
func GetSpecificRegistration(ID, UUID string) (*structs.CountryInfoInternal, error) {
	ref := Client.Collection(Firestore.RegistrationCollection) // Reference to the Registration collection.

	// Query for the specific document with the given 'ID' and 'UUID'.
	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	var ci *structs.CountryInfoInternal // Variable to store the fetched document.

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return nil, fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		if err := doc.DataTo(&ci); err != nil {
			return nil, err // Return error if parsing the document fails.
		}
		log.Printf("Registration document %s retrieved successfully.", ci.ID)
		return ci, nil // Return the parsed document.
	}

	return nil, errors.New("no registration with that ID was found") // Return error if no document is found.
}

// UpdateRegistration updates a specific registration document by ID and UUID in Firestore.
func UpdateRegistration(ID, UUID string, data *structs.CountryInfoInternal) error {
	ref := Client.Collection(Firestore.RegistrationCollection) // Reference to the Registration collection.

	// Query for the specific document to update.
	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		// Update the document with the provided data.
		_, err = ref.Doc(doc.Ref.ID).Set(ctx, map[string]interface{}{
			"ID":         data.ID,
			"UUID":       data.UUID,
			"Country":    data.Country,
			"IsoCode":    data.IsoCode,
			"Features":   data.Features,
			"Lastchange": firestore.ServerTimestamp, // Use server timestamp to update 'Lastchange'.
		})
		if err != nil {
			return err // Return error if the document update operation fails.
		}
		log.Printf("Registration document %s patched successfully.", doc.Ref.ID)
		return nil // Return nil error if the update is successful.
	}

	return errors.New("no registration with that ID was found") // Return error if no document is found.
}

// DeleteRegistration deletes a specific registration document by ID and UUID from Firestore.
func DeleteRegistration(ID, UUID string) error {
	ref := Client.Collection(Firestore.RegistrationCollection) // Reference to the Registration collection.

	// Query for the specific document to delete.
	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	var docID string // Variable to store the document ID of the found registration.
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		docID = doc.Ref.ID // Store the document ID.
	}

	if docID == "" {
		return fmt.Errorf("ID match was not found") // Return error if no document ID was found.
	}

	_, err = ref.Doc(docID).Delete(ctx) // Delete the document from Firestore.
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err) // Return formatted error if delete fails.
	}

	log.Printf("Registration document %s deleted successfully\n", docID)
	return nil // Return nil if the deletion is successful.
}

// AddWebhook creates a new webhook entry in Firestore.
func AddWebhook(docID string, webhook *structs.WebhookInternal) error {
	ref := Client.Collection(Firestore.WebhookCollection) // Reference to the Webhook collection in Firestore.

	// Set the webhook document with the provided ID and data.
	_, err = ref.Doc(docID).Set(ctx, webhook)
	if err != nil {
		return err // Return error if addition fails.
	}

	log.Printf("Webhook %s created successfully.", webhook.ID) // Log success.
	return nil                                                 // Return nil error on successful addition.
}

// GetAllWebhooks retrieves all webhook entries from Firestore.
func GetAllWebhooks() ([]structs.WebhookInternal, error) {
	ref := Client.Collection(Firestore.WebhookCollection) // Reference to the Webhook collection.

	// Retrieve all documents from the webhook collection.
	docs, err := ref.Documents(ctx).GetAll()
	if err != nil {
		return nil, err // Return the error if the fetch operation fails.
	}

	var webhooks []structs.WebhookInternal // Slice to store the fetched webhook documents.

	for _, doc := range docs {
		var webhook structs.WebhookInternal
		if err := doc.DataTo(&webhook); err != nil {
			return nil, err // Return error if parsing any document fails.
		}
		webhooks = append(webhooks, webhook) // Append the parsed document to the slice.
	}

	log.Printf("All Webhooks retrieved successfully.") // Log success.
	return webhooks, nil                               // Return the slice of webhook documents.
}

// GetWebhooksUser retrieves all webhook entries for a specific user (UUID) from Firestore.
func GetWebhooksUser(UUID string) ([]structs.WebhookResponse, error) {
	ref := Client.Collection(Firestore.WebhookCollection) // Reference to the Webhook collection.

	// Query and retrieve all documents from the webhook collection where 'UUID' matches the provided UUID.
	docs, err := ref.Where("UUID", "==", UUID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err // Return the error if the fetch operation fails.
	}

	var webhooks []structs.WebhookResponse // Slice to store the fetched webhook documents for the user.

	for _, doc := range docs {
		var webhook structs.WebhookResponse
		if err := doc.DataTo(&webhook); err != nil {
			return nil, err // Return error if parsing any document fails.
		}
		webhooks = append(webhooks, webhook) // Append the parsed document to the slice.
	}

	log.Printf("Webhooks retrieved successfully for user: %s.", UUID) // Log success.
	return webhooks, nil                                              // Return the slice of webhook documents for the user.
}

// GetSpecificWebhook retrieves a specific webhook entry by ID and UUID from Firestore.
func GetSpecificWebhook(ID, UUID string) (*structs.WebhookResponse, error) {
	ref := Client.Collection(Firestore.WebhookCollection) // Reference to the Webhook collection.

	// Query for the specific document with the given 'ID' and 'UUID'.
	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	var webhook *structs.WebhookResponse // Variable to store the fetched document.

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return nil, fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		if err := doc.DataTo(&webhook); err != nil {
			return nil, err // Return error if parsing the document fails.
		}

		log.Printf("Webhook %s retrieved successfully.", webhook.ID) // Log success.
		return webhook, nil                                          // Return the parsed document.
	}

	return nil, errors.New("no document with that ID was found") // Return error if no document is found.
}

// DeleteWebhook deletes a specific webhook entry by ID and UUID from Firestore.
func DeleteWebhook(ID, UUID string) error {
	ref := Client.Collection(Firestore.WebhookCollection) // Reference to the Webhook collection.

	// Query for the specific document to delete.
	iter := ref.Where("ID", "==", ID).Where("UUID", "==", UUID).Limit(1).Documents(ctx)
	defer iter.Stop() // Ensure the iterator is cleaned up properly.

	var docID string // Variable to store the document ID of the found webhook.
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop if all documents have been iterated over.
		}
		if err != nil {
			return fmt.Errorf(IterationFailed, err) // Return formatted error if iteration fails.
		}
		docID = doc.Ref.ID // Store the document ID.
	}

	if docID == "" {
		return fmt.Errorf("ID match was not found") // Return error if no document ID was found.
	}

	_, err = ref.Doc(docID).Delete(ctx) // Delete the document from Firestore.
	if err != nil {
		log.Print(err)                                          // Log errors during the delete operation.
		return fmt.Errorf("failed to delete document: %v", err) // Return formatted error if delete fails.
	}

	log.Printf("Webhook %s deleted successfully.", docID) // Log success.
	return nil                                            // Return nil error on successful operation.
}
