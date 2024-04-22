// Package dashboard provides handlers for managing dashboard-related functionalities through HTTP endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Webhooks"
	"globeboard/internal/utils/structs"
	"io"
	"log"
	"net/http"
)

// Constant strings used for API responses and header configurations
const (
	ProvideAPI      = "Please provide API Token" // Message prompting for API token
	APINotAccepted  = "API key not accepted"     // Message for an unauthorized API token
	ContentType     = "Content-Type"             // HTTP header field for content-type
	ApplicationJSON = "application/json"         // MIME type for JSON
)

// RegistrationsHandler routes the HTTP request based on the method (POST, GET) to appropriate handlers
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleRegPostRequest(w, r) // Handle POST requests
	case http.MethodGet:
		handleRegGetAllRequest(w, r) // Handle GET requests
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, Endpoints.Registrations, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet+"\n"+http.MethodPatch, http.StatusNotImplemented)
		return
	}
}

// DecodeCountryInfo decodes JSON data from the request body into a CountryInfoInternal struct
func DecodeCountryInfo(data io.ReadCloser) (*structs.CountryInfoInternal, error) {
	var ci *structs.CountryInfoInternal
	if err := json.NewDecoder(data).Decode(&ci); err != nil {
		return nil, err // Return error if decoding fails
	}

	err := _func.ValidateCountryInfo(ci) // Validate the decoded information
	if err != nil {
		return nil, err // Return validation errors
	}

	return ci, nil // Return the decoded and validated country information
}

// handleRegPostRequest handles the POST requests for registration endpoint
func handleRegPostRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Extract the 'token' parameter from the query.
	if token == "" {            // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.Method, Endpoints.Registrations)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token) // Retrieve the UUID for the API key.
	if UUID == "" {                 // Validate UUID presence.
		log.Printf(constants.ClientConnectUnauthorized, r.Method, Endpoints.Registrations)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}

	if r.Body == nil { // Validate that the request body is not empty.
		log.Printf(constants.ClientConnectEmptyBody, r.Method, Endpoints.Registrations)
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	ci, err := DecodeCountryInfo(r.Body) // Decode request body into CountryInfoInternal struct.
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		err := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	UDID := _func.GenerateUID(constants.DocIdLength) // Generate a unique ID for document
	URID := _func.GenerateUID(constants.IdLength)    // Generate a unique ID for registration

	ci.ID = URID
	ci.UUID = UUID

	err = db.AddRegistration(UDID, ci) // Add Registration to the Database.
	if err != nil {
		log.Println("Error saving data to database" + err.Error())
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	reg, err := db.GetSpecificRegistration(URID, UUID) // Retrieve specific registration details.
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error confirming data added to database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{ // construct JSON response.
		"id":         reg.ID,
		"lastChange": reg.Lastchange,
	}

	w.Header().Set(ContentType, ApplicationJSON) // Set the content type of response

	w.WriteHeader(http.StatusCreated) // Set HTTP status code to 201 Created

	err = json.NewEncoder(w).Encode(response) // Encode and send the response
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cie := new(structs.CountryInfoExternal) // Create new external country info struct.
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.Registrations, Webhooks.EventRegister) // Send webhook notifications
}

// handleRegGetAllRequest handles the GET requests for registration endpoint to retrieve all registrations
func handleRegGetAllRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Extract the 'token' parameter from the query.
	if token == "" {            // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.Method, Endpoints.Registrations)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token) // Extract the UUID parameter from the API token.
	if UUID == "" {                 // Validate UUID presence.
		log.Printf(constants.ClientConnectUnauthorized, r.Method, Endpoints.Registrations)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	regs, err := db.GetRegistrations(UUID) // Retrieve the user's Registrations.
	if err != nil {
		log.Printf("Error retrieving documents from database: %s", err)
		errmsg := fmt.Sprint("Error retrieving documents from database: ", err)
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON) // Set the content type of response

	w.WriteHeader(http.StatusOK) // Set HTTP status code to 200 OK

	var cies []*structs.CountryInfoExternal // Construct CountryInfoExternal slice.
	for _, reg := range regs {              // Loop over the retrieved registrations and parse them to CountryInfoExternal struct.
		cie := new(structs.CountryInfoExternal)
		cie.ID = reg.ID
		cie.Country = reg.Country
		cie.IsoCode = reg.IsoCode
		cie.Features = reg.Features
		cie.Lastchange = reg.Lastchange
		cies = append(cies, cie) // Append the individual CountryInfoExternal structs to the slice.
	}

	err = json.NewEncoder(w).Encode(cies) // Encode and send the response
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, cie := range cies {
		_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.Registrations, Webhooks.EventInvoke) // Send webhook notifications on data retrieval
	}
}
