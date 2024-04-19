// Package dashboard provides handlers for dashboard-related endpoints.
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

const (
	ProvideAPI      = "Please provide API Token"
	APINotAccepted  = "API key not accepted"
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

// RegistrationsHandler handles HTTP GET requests to retrieve supported languages.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleRegPostRequest(w, r)
	case http.MethodGet:
		handleRegGetAllRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet+"\n"+http.MethodPatch, http.StatusNotImplemented)
		return
	}
}

func DecodeCountryInfo(data io.ReadCloser) (*structs.CountryInfoInternal, error) {
	var ci *structs.CountryInfoInternal
	if err := json.NewDecoder(data).Decode(&ci); err != nil {
		return nil, err
	}

	err := _func.ValidateCountryInfo(ci)
	if err != nil {
		return nil, err
	}

	return ci, nil
}

// handleRegPostRequest handles POST requests to register a country.
func handleRegPostRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token)
	if UUID == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}

	if r.Body == nil {
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// Decode the body into a CountryInfoPost struct
	ci, err := DecodeCountryInfo(r.Body)
	if err != nil {
		err := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	UDID := _func.GenerateUID(constants.DocIdLength)
	URID := _func.GenerateUID(constants.IdLength)

	ci.ID = URID
	ci.UUID = UUID

	err = db.AddRegistration(UDID, ci)
	if err != nil {
		log.Println("Error saving data to database" + err.Error())
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	reg, err := db.GetSpecificRegistration(URID, UUID)
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error confirming data added to database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         reg.ID,
		"lastChange": reg.Lastchange,
	}

	// Set Content-Type header
	w.Header().Set(ContentType, ApplicationJSON)

	// Write the status code to the response
	w.WriteHeader(http.StatusCreated)

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cie := new(structs.CountryInfoExternal)
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.Registrations, Webhooks.EventRegister)
}

// handleRegGetAllRequest handles GET requests to retrieve a registered country.
func handleRegGetAllRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token)
	if UUID == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	regs, err := db.GetRegistrations(UUID)
	if err != nil {
		errmsg := fmt.Sprint("Error retrieving document from database: ", err)
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set(ContentType, ApplicationJSON)

	// Write the status code to the response
	w.WriteHeader(http.StatusOK)

	// Parse Data for External Users
	var cies []*structs.CountryInfoExternal
	for _, reg := range regs {
		cie := new(structs.CountryInfoExternal)
		cie.ID = reg.ID
		cie.Country = reg.Country
		cie.IsoCode = reg.IsoCode
		cie.Features = reg.Features
		cie.Lastchange = reg.Lastchange
		cies = append(cies, cie)
	}

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(cies)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, cie := range cies {
		_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.Registrations, Webhooks.EventInvoke)
	}
}
