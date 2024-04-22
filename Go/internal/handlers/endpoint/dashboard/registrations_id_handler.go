// Package dashboard provides handlers for managing dashboard-related functionalities through HTTP endpoints.
package dashboard

import (
	"encoding/json"
	"errors"
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

// Constants for error and informational messages.
const (
	ProvideID                 = "Please Provide ID"                   // Message to request ID provision when missing.
	RegistrationRetrivalError = "%s: Error getting registration: %v"  // Error message template for retrieval issues.
	RegistrationPatchError    = "%s: Error patching registration: %v" // Error message template for patching issues.
)

// RegistrationsIdHandler handles requests for the /registrations/{ID} endpoint.
func RegistrationsIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // Handle GET requests.
		handleRegGetRequest(w, r)
	case http.MethodPatch: // Handle PATCH requests.
		handleRegPatchRequest(w, r)
	case http.MethodDelete: // Handle DELETE requests.
		handleRegDeleteRequest(w, r)
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.RegistrationsID, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodGet+"\n"+http.MethodPatch+"\n"+http.MethodDelete, http.StatusNotImplemented)
		return
	}
}

// handleRegGetRequest processes GET requests for registration data by ID.
func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")     // Extract the 'ID' parameter from the URL path.
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Extract the 'token' parameter from the query.
	if token == "" {            // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(r.RemoteAddr, token) // Retrieve the UUID for the API key.
	if UUID == "" {                               // Validate UUID presence.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" || ID == " " { // Validate ID presence.
		log.Printf(constants.ClientConnectNoID, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(r.RemoteAddr, ID, UUID) // Retrieve registration data from the database.
	if err != nil {
		log.Printf(RegistrationRetrivalError, r.RemoteAddr, err)
		http.Error(w, "Error retrieving data from database", http.StatusNotFound)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON) // Set the Content-Type header.

	w.WriteHeader(http.StatusOK) // Set the HTTP status code to 200.

	cie := new(structs.CountryInfoExternal) // Create new external country info struct.
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features
	cie.Lastchange = reg.Lastchange

	err = json.NewEncoder(w).Encode(cie) // Encode the external country info into JSON and write to the response.
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.RegistrationsID, Webhooks.EventInvoke) // Trigger webhooks for the registration.
}

// handleRegPatchRequest processes PATCH requests to update registration data by ID.
func handleRegPatchRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")     // Extract the 'ID' parameter from the URL path.
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Extract the 'token' parameter from the query.
	if token == "" {            // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(r.RemoteAddr, token) // Retrieve the UUID for the API key.
	if UUID == "" {                               // Validate UUID presence.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" || ID == " " { // Validate ID presence.
		log.Printf(constants.ClientConnectNoID, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	if r.Body == nil { // Validate that the request body is not empty.
		log.Printf(constants.ClientConnectEmptyBody, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	ci, err, errcode := patchCountryInformation(r, ID, UUID) // Process the patch request.
	if err != nil {
		log.Printf(RegistrationPatchError, r.RemoteAddr, err)
		err := fmt.Sprintf("Error patching registration: %v", err)
		http.Error(w, err, errcode)
		return
	}

	err = _func.ValidateCountryInfo(ci) // Validate the patched country information.
	if err != nil {
		log.Printf(RegistrationPatchError, r.RemoteAddr, err)
		err := fmt.Sprintf("Error patching registration: %v", err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	err = db.UpdateRegistration(r.RemoteAddr, ID, UUID, ci) // Update the registration in the database.
	if err != nil {
		log.Printf("%s: Error saving patched data to database: %v", r.RemoteAddr, err)
		err := fmt.Sprintf("Error saving patched data to database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	reg, err := db.GetSpecificRegistration(r.RemoteAddr, ID, UUID) // Retrieve the updated registration.
	if err != nil {
		log.Printf(RegistrationRetrivalError, r.RemoteAddr, err)
		err := fmt.Sprint("Error retrieving updated document: ", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	cie := new(structs.CountryInfoExternal) // Create new external country info struct.
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features
	cie.Lastchange = reg.Lastchange

	w.Header().Set("content-type", "application/json") // Set the Content-Type header.
	w.WriteHeader(http.StatusAccepted)                 // Set the HTTP status code to 202.

	response := map[string]interface{}{
		"lastChange": cie.Lastchange, // Prepare the response data.
	}

	err = json.NewEncoder(w).Encode(response) // Encode the response data into JSON and write to the response.
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.RegistrationsID, Webhooks.EventChange) // Trigger webhooks for the change event.
}

// patchCountryInformation updates the country information based on the provided patch data.
func patchCountryInformation(r *http.Request, ID, UUID string) (*structs.CountryInfoInternal, error, int) {
	reg, err := db.GetSpecificRegistration(r.RemoteAddr, ID, UUID) // Retrieve the specific registration.
	if err != nil {
		log.Printf(RegistrationRetrivalError, r.RemoteAddr, err)
		return nil, err, http.StatusNotFound
	}

	bytes, err := json.Marshal(reg) // Marshal the registration data to JSON bytes.
	if err != nil {
		log.Print(err)
		return nil, err, http.StatusInternalServerError
	}

	var originalData map[string]interface{} // Unmarshal the JSON bytes back to a map.
	err = json.Unmarshal(bytes, &originalData)
	if err != nil {
		log.Print(err)
		return nil, err, http.StatusInternalServerError
	}

	all, err := io.ReadAll(r.Body) // Read all data from the request body.
	if err != nil {
		log.Print(err)
		return nil, err, http.StatusInternalServerError
	}

	var patchData map[string]interface{} // Unmarshal the patch data from the request body.
	err = json.Unmarshal(all, &patchData)
	if err != nil {
		log.Print(err)
		return nil, err, http.StatusInternalServerError
	}

	patchFeatures, err, errcode := validatePatchData(patchData, originalData) // Validate and extract the patch data.
	if err != nil {
		return nil, err, errcode
	}

	if originalData["features"] != nil { // Merge the patch features into the original features.
		originalFeatures := originalData["features"].(map[string]interface{})
		for key, value := range patchFeatures {
			originalFeatures[key] = value
		}
	}

	jsonData, err := json.Marshal(originalData) // Marshal the updated data to JSON.
	if err != nil {
		log.Print(err)
		return nil, err, http.StatusInternalServerError
	}

	var countryInfo *structs.CountryInfoInternal // Unmarshal the JSON data to a CountryInfoInternal struct.
	err = json.Unmarshal(jsonData, &countryInfo)
	if err != nil {
		log.Print(err)
		return nil, err, http.StatusInternalServerError
	}
	return countryInfo, nil, http.StatusOK
}

// validatePatchData checks the validity of the patch data against the original data.
func validatePatchData(patchData map[string]interface{}, originalData map[string]interface{}) (map[string]interface{}, error, int) {
	// Check if "country" or "isoCode" fields are provided and if they are non-empty and differ from the original data.
	if country, ok := patchData["country"]; ok {
		if countryStr, isStr := country.(string); isStr && countryStr != "" && originalData["country"] != country {
			return nil, errors.New("modification of 'country' field is not allowed"), http.StatusBadRequest
		}
	}

	if isoCode, ok := patchData["isocode"]; ok {
		if isoCodeStr, isStr := isoCode.(string); isStr && isoCodeStr != "" && originalData["isocode"] != isoCode {
			return nil, errors.New("modification of 'isocode' field is not allowed"), http.StatusBadRequest
		}
	}

	// Enforce "features" to be provided and not empty.
	features, ok := patchData["features"]
	if !ok || features == nil {
		return nil, errors.New("user must provide features to patch"), http.StatusBadRequest
	}
	patchFeatures, isMap := features.(map[string]interface{})
	if !isMap || len(patchFeatures) == 0 {
		return nil, errors.New("user must provide non-empty features to patch"), http.StatusBadRequest
	}
	return patchFeatures, nil, http.StatusOK
}

// handleRegDeleteRequest processes DELETE requests to remove registration data by ID.
func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")     // Extract the 'ID' parameter from the URL path.
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Extract the 'token' parameter from the query.
	if token == "" {            // Validate token presence.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(r.RemoteAddr, token) // Retrieve the UUID for the API key.
	if UUID == "" {                               // Validate UUID presence.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" || ID == " " { // Validate ID presence.
		log.Printf(constants.ClientConnectNoID, r.RemoteAddr, r.Method, Endpoints.RegistrationsID)
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(r.RemoteAddr, ID, UUID) // Retrieve the specific registration to be deleted.
	if err != nil {
		log.Printf(RegistrationRetrivalError, r.RemoteAddr, err)
		err := fmt.Sprint("Error getting registration: ", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	err = db.DeleteRegistration(r.RemoteAddr, ID, UUID) // Delete the registration from the database.
	if err != nil {
		log.Printf("%s: Error deleting registration from database: %v", r.RemoteAddr, err)
		err := fmt.Sprintf("Error deleting registration from database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Set the HTTP status code to 204 (No Content).

	cie := new(structs.CountryInfoExternal) // Create new external country info struct.
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features
	cie.Lastchange = reg.Lastchange

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.RegistrationsID, Webhooks.EventDelete) // Trigger webhooks for the delete event.
}
