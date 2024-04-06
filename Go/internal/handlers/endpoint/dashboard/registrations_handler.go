// Package dashboard provides handlers for dashboard-related endpoints.
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
	"net/http"
	"time"
)

// RegistrationsHandler handles HTTP GET requests to retrieve supported languages.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleRegPostRequest(w, r)
	case http.MethodGet:
		handleRegGetRequest(w, r)
	case http.MethodPut:
		handleRegPutRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported method for this endpoint is:\n"+http.MethodPost+"\n"+http.MethodGet+"\n"+http.MethodPut, http.StatusNotImplemented)
		return
	}
}

func DecodeCountryInfo(data io.ReadCloser) (*structs.CountryInfoPost, error) {
	var ci *structs.CountryInfoPost
	if err := json.NewDecoder(data).Decode(&ci); err != nil {
		return nil, err
	}

	err := validateCountryNameIsoCode(ci)
	if err != nil {
		return nil, err
	}

	if !ci.Features.Temperature && !ci.Features.Precipitation && !ci.Features.Capital &&
		!ci.Features.Coordinates && !ci.Features.Population && !ci.Features.Area {
		return nil, errors.New("at least one feature must be true")
	}

	return ci, nil
}

func validateCountryNameIsoCode(ci *structs.CountryInfoPost) error {
	validCountries, err := _func.GetSupportedCountries() // Adjusted to use the map version.
	if err != nil {
		return errors.New("error validating country")
	}

	if err := validateCountryOrIsoCodeProvided(ci); err != nil {
		return err
	}

	if err := validateIsoCode(ci, validCountries); err != nil {
		return err
	}

	if err := updateAndValidateIsoCodeForCountry(ci, validCountries); err != nil {
		return err
	}

	return validateCorrespondence(ci, validCountries)
}

func validateCountryOrIsoCodeProvided(ci *structs.CountryInfoPost) error {
	if ci.Country == "" && ci.IsoCode == "" {
		return errors.New("either country name or ISO code must be provided")
	}
	return nil
}

func validateIsoCode(ci *structs.CountryInfoPost, validCountries map[string]string) error {
	if ci.IsoCode != "" {
		if _, exists := validCountries[ci.IsoCode]; !exists {
			return errors.New("invalid ISO code")
		}
	}
	return nil
}

func updateAndValidateIsoCodeForCountry(ci *structs.CountryInfoPost, validCountries map[string]string) error {
	if ci.IsoCode == "" && ci.Country != "" {
		for code, name := range validCountries {
			if name == ci.Country {
				ci.IsoCode = code
				return nil
			}
		}
		return errors.New("country name not valid or not supported")
	}
	return nil
}

func validateCorrespondence(ci *structs.CountryInfoPost, validCountries map[string]string) error {
	if ci.Country != "" && ci.IsoCode != "" {
		if validCountries[ci.IsoCode] != ci.Country {
			return errors.New("ISO code and country name do not match")
		}
	}
	return nil
}

// handleRegPostRequest handles POST requests to register a country.
func handleRegPostRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, "Please provide API Token", http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf("Error checking API key: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !exists {
		err := fmt.Sprintf("API key not accepted")
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

	err = db.AddRegistration(UDID, UUID, ci)
	if err != nil {
		http.Error(w, "Error storing data in database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         UUID,
		"lastChange": time.Now(),
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response
	w.WriteHeader(http.StatusCreated)

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooks(ci, Endpoints.Registrations, Webhooks.EventRegister)
}

// handleRegGetRequest handles GET requests to retrieve a registered country.
func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	//TODO::Complete HTTP Method Requests
}

// handleRegPutRequest handles PUT requests to Update a registered country.
func handleRegPutRequest(w http.ResponseWriter, r *http.Request) {
	//TODO::Complete HTTP Method Requests
}
