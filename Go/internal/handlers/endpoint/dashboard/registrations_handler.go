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
	"net/http"
	"strings"
	"time"
)

const (
	ProvideAPI        = "Please provide API Token"
	ErrorAPI          = "Error checking API key: %v"
	APINotAccepted    = "API key not accepted"
	ContentType       = "Content-Type"
	ApplicationJSON   = "application/json"
	InvalidURL        = "Invalid URL"
	DataRetrivalError = "Error retrieving data from database"
)

// RegistrationsHandler handles HTTP GET requests to retrieve supported languages.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleRegPostRequest(w, r)
	case http.MethodGet:
		handleRegGetRequest(w, r)
	case http.MethodPatch:
		handleRegPatchRequest(w, r)
	case http.MethodDelete:
		handleRegDeleteRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet+"\n"+http.MethodPatch, http.StatusNotImplemented)
		return
	}
}

func DecodeCountryInfo(data io.ReadCloser) (*structs.CountryInfoGet, error) {
	var ci *structs.CountryInfoGet
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
		http.Error(w, ProvideAPI, http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf(ErrorAPI, err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !exists {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}

	if r.Body == nil {
		err := fmt.Sprintf("Please send a request body: %v", err)
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
	UUID := _func.GenerateUID(constants.RegIdLength)

	lastchange := time.Now()

	ci.Id = UUID
	ci.Lastchange = lastchange

	err = db.AddRegistration(UDID, ci)
	if err != nil {
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         UUID,
		"lastChange": lastchange,
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

	_func.LoopSendWebhooks(ci, Endpoints.RegistrationsSlash, Webhooks.EventRegister)
}

// handleRegGetRequest handles GET requests to retrieve a registered country.
func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	registrationId := ""
	if len(parts) >= 5 {
		registrationId = parts[4] // Language code will be at index 4
	} else {
		http.Error(w, InvalidURL, http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf(ErrorAPI, err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !exists {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if registrationId == "" {
		GetAllRegistrations(w)
	} else {
		GetSpecificRegistrations(w, registrationId)
	}
}

func GetAllRegistrations(w http.ResponseWriter) {
	regs, err := db.GetRegistrations()
	if err != nil {
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set(ContentType, ApplicationJSON)

	// Write the status code to the response
	w.WriteHeader(http.StatusOK)

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(regs)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, reg := range regs {
		_func.LoopSendWebhooks(reg, Endpoints.RegistrationsSlash, Webhooks.EventInvoke)
	}
}

func GetSpecificRegistrations(w http.ResponseWriter, id string) {
	reg, err := db.GetSpecificRegistration(id)
	if err != nil {
		http.Error(w, DataRetrivalError, http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set(ContentType, ApplicationJSON)

	// Write the status code to the response
	w.WriteHeader(http.StatusOK)

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(reg)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooks(reg, Endpoints.RegistrationsSlash, Webhooks.EventInvoke)
}

// handleRegPatchRequest handles PUT requests to Update a registered country.
func handleRegPatchRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	registrationId := ""
	if len(parts) >= 5 {
		registrationId = parts[4] // Language code will be at index 4
	} else {
		http.Error(w, InvalidURL, http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf(ErrorAPI, err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !exists {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if registrationId == "" {
		http.Error(w, "Please provide ID", http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		err := fmt.Sprintf("Please send a request body: %v", err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	countryInfo, err := patchCountryInformation(r, err, registrationId)
	if err != nil {
		err := fmt.Sprintf("Error patching data together: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	err = _func.ValidateCountryInfo(countryInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	countryInfo.Lastchange = time.Now()

	err = db.UpdateRegistration(registrationId, countryInfo)
	if err != nil {
		err := fmt.Sprintf("Error saving patched data to database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	_func.LoopSendWebhooks(countryInfo, Endpoints.RegistrationsSlash, Webhooks.EventChange)
}

func patchCountryInformation(r *http.Request, err error, registrationId string) (*structs.CountryInfoGet, error) {
	reg, err := db.GetSpecificRegistration(registrationId)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}

	var originalData map[string]interface{}
	err = json.Unmarshal(bytes, &originalData)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var patchData map[string]interface{}
	err = json.Unmarshal(all, &patchData)
	if err != nil {
		return nil, err
	}

	for key, value := range patchData {
		// If the key is for a nested object, handle it appropriately
		if key == "features" && originalData[key] != nil {
			// Assume both are maps, merge them
			for subKey, subValue := range value.(map[string]interface{}) {
				originalData[key].(map[string]interface{})[subKey] = subValue
			}
		} else {
			originalData[key] = value
		}
	}

	// First, marshal your map into JSON
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		return nil, err
	}

	var countryInfo *structs.CountryInfoGet
	err = json.Unmarshal(jsonData, &countryInfo)
	if err != nil {
		return nil, err
	}
	return countryInfo, nil
}

func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	registrationId := ""
	if len(parts) >= 5 {
		registrationId = parts[4] // Language code will be at index 4
	} else {
		http.Error(w, InvalidURL, http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf(ErrorAPI, err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !exists {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if registrationId == "" {
		http.Error(w, "Please provide ID", http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(registrationId)
	if err != nil {
		http.Error(w, DataRetrivalError, http.StatusInternalServerError)
		return
	}

	err = db.DeleteRegistration(registrationId)
	if err != nil {
		err := fmt.Sprintf("Error saving patched data to database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	_func.LoopSendWebhooks(reg, Endpoints.RegistrationsSlash, Webhooks.EventDelete)

}
