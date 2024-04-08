package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Webhooks"
	"globeboard/internal/utils/structs"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	ProvideID = "Please Provide ID"
)

// RegistrationsIdHandler handles HTTP GET requests to retrieve supported languages.
func RegistrationsIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRegGetRequest(w, r)
	case http.MethodPatch:
		handleRegPatchRequest(w, r)
	case http.MethodDelete:
		handleRegDeleteRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodGet+"\n"+http.MethodPatch+"\n"+http.MethodDelete, http.StatusNotImplemented)
		return
	}
}

// handleRegGetRequest handles GET requests to retrieve a registered country.
func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	uuid := db.GetAPIKeyUUID(token)
	if uuid == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if id == "" {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(id, uuid)
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error retrieving data from database", http.StatusInternalServerError)
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

	_func.LoopSendWebhooks(uuid, reg, Endpoints.RegistrationsID, Webhooks.EventInvoke)
}

// handleRegPatchRequest handles PUT requests to Update a registered country.
func handleRegPatchRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	uuid := db.GetAPIKeyUUID(token)
	if uuid == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if id == "" {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	countryInfo, err := patchCountryInformation(r, id, uuid)
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

	err = db.UpdateRegistration(id, uuid, countryInfo)
	if err != nil {
		err := fmt.Sprintf("Error saving patched data to database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	_func.LoopSendWebhooks(uuid, countryInfo, Endpoints.RegistrationsID, Webhooks.EventChange)
}

func patchCountryInformation(r *http.Request, registrationId, uuid string) (*structs.CountryInfoGet, error) {
	reg, err := db.GetSpecificRegistration(registrationId, uuid)
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

	patchFeatures, err := validatePatchData(patchData, originalData)
	if err != nil {
		return nil, err
	}

	if originalData["features"] != nil {
		originalFeatures := originalData["features"].(map[string]interface{})
		for key, value := range patchFeatures {
			originalFeatures[key] = value
		}
	}

	// Marshal the original data back to JSON.
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the CountryInfoGet struct.
	var countryInfo *structs.CountryInfoGet
	err = json.Unmarshal(jsonData, &countryInfo)
	if err != nil {
		return nil, err
	}
	return countryInfo, nil
}

func validatePatchData(patchData map[string]interface{}, originalData map[string]interface{}) (map[string]interface{}, error) {
	// Check if "country" or "isoCode" fields are provided and if they are non-empty and differ from the original data.
	if country, ok := patchData["country"]; ok {
		if countryStr, isStr := country.(string); isStr && countryStr != "" && originalData["country"] != country {
			return nil, errors.New("modification of 'country' field is not allowed")
		}
	}

	if isoCode, ok := patchData["isoCode"]; ok {
		if isoCodeStr, isStr := isoCode.(string); isStr && isoCodeStr != "" && originalData["isoCode"] != isoCode {
			return nil, errors.New("modification of 'isoCode' field is not allowed")
		}
	}

	// Enforce "features" to be provided and not empty.
	features, ok := patchData["features"]
	if !ok || features == nil {
		return nil, errors.New("user must provide features to patch")
	}
	patchFeatures, isMap := features.(map[string]interface{})
	if !isMap || len(patchFeatures) == 0 {
		return nil, errors.New("user must provide non-empty features to patch")
	}
	return patchFeatures, nil
}

func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	uuid := db.GetAPIKeyUUID(token)
	if uuid == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if id == "" {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(id, uuid)
	if err != nil {
		log.Println("Document doesn't exist: ", err)
		http.Error(w, "", http.StatusNoContent)
		return
	}

	err = db.DeleteRegistration(id, uuid)
	if err != nil {
		err := fmt.Sprintf("Error saving patched data to database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	_func.LoopSendWebhooks(uuid, reg, Endpoints.RegistrationsID, Webhooks.EventDelete)
}
