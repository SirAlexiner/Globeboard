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
)

const (
	ProvideID = "Please Provide ID"
)

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

func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
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
	if ID == "" || ID == " " {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(ID, UUID)
	if err != nil {
		log.Print("Error getting document from database: ", err)
		http.Error(w, "Error retrieving data from database", http.StatusNotFound)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON)

	w.WriteHeader(http.StatusOK)

	cie := new(structs.CountryInfoExternal)
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features
	cie.Lastchange = reg.Lastchange

	err = json.NewEncoder(w).Encode(cie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.RegistrationsID, Webhooks.EventInvoke)
}

func handleRegPatchRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
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
	if ID == "" || ID == " " {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		err := fmt.Sprintf("Please send a request body")
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	ci, err, errcode := patchCountryInformation(r, ID, UUID)
	if err != nil {
		err := fmt.Sprintf("Error patching data together: %v", err)
		http.Error(w, err, errcode)
		return
	}

	err = _func.ValidateCountryInfo(ci)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateRegistration(ID, UUID, ci)
	if err != nil {
		err := fmt.Sprintf("Error saving patched data to database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	reg, err := db.GetSpecificRegistration(ID, UUID)
	if err != nil {
		err := fmt.Sprint("Error retrieving updated document: ", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	cie := new(structs.CountryInfoExternal)
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features
	cie.Lastchange = reg.Lastchange

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]interface{}{
		"lastChange": cie.Lastchange,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.RegistrationsID, Webhooks.EventChange)
}

func patchCountryInformation(r *http.Request, ID, UUID string) (*structs.CountryInfoInternal, error, int) {
	reg, err := db.GetSpecificRegistration(ID, UUID)
	if err != nil {
		return nil, err, http.StatusNotFound
	}

	bytes, err := json.Marshal(reg)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var originalData map[string]interface{}
	err = json.Unmarshal(bytes, &originalData)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	all, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var patchData map[string]interface{}
	err = json.Unmarshal(all, &patchData)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	patchFeatures, err, errcode := validatePatchData(patchData, originalData)
	if err != nil {
		return nil, err, errcode
	}

	if originalData["features"] != nil {
		originalFeatures := originalData["features"].(map[string]interface{})
		for key, value := range patchFeatures {
			originalFeatures[key] = value
		}
	}

	jsonData, err := json.Marshal(originalData)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var countryInfo *structs.CountryInfoInternal
	err = json.Unmarshal(jsonData, &countryInfo)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return countryInfo, nil, http.StatusOK
}

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

func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
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
	if ID == "" || ID == " " {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(ID, UUID)
	if err != nil {
		err := fmt.Sprint("Document doesn't exist: ", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	err = db.DeleteRegistration(ID, UUID)
	if err != nil {
		err := fmt.Sprintf("Error deleting data from database: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	cie := new(structs.CountryInfoExternal)
	cie.ID = reg.ID
	cie.Country = reg.Country
	cie.IsoCode = reg.IsoCode
	cie.Features = reg.Features
	cie.Lastchange = reg.Lastchange

	_func.LoopSendWebhooksRegistrations(UUID, cie, Endpoints.RegistrationsID, Webhooks.EventDelete)
}
