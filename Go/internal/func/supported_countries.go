// Package _func provides developer-made utility functions for use within the application.
package _func

import (
	"encoding/json"
	"errors"
	"fmt"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// getSupportedCountries fetches supported countries with their common names and ISO 3166-1 alpha-2 codes.
func getSupportedCountries() (map[string]string, error) {
	url := fmt.Sprintf("%sall?fields=name,cca2", External.CountriesAPI) // Constructing the API request URL.
	var responseData []struct {                                         // Struct to parse the JSON response.
		Name struct {
			Common string `json:"common"` // Common name of the country.
		} `json:"name"`
		CCA2 string `json:"cca2"` // ISO 3166-1 alpha-2 code of the country.
	}

	client := &http.Client{Timeout: 10 * time.Second}     // HTTP client with a timeout.
	req, err := http.NewRequest(http.MethodGet, url, nil) // Creating a new HTTP GET request.
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("content-type", "application/json") // Setting content-type of the request.

	res, err := client.Do(req) // Send the request.
	if err != nil {
		log.Printf("Error issuing request: %v", err)
		return nil, fmt.Errorf("error issuing request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(res.Body) // Ensure the response body is closed after processing.

	err = json.NewDecoder(res.Body).Decode(&responseData) // Decoding the JSON response into the struct.
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	countriesMap := make(map[string]string) // Map to hold country codes and their common names.
	// Loop over the response data and map the supported countries by ISO code.
	for _, item := range responseData {
		countriesMap[item.CCA2] = item.Name.Common
	}
	return countriesMap, nil // Returning the map.
}

// ValidateCountryInfo validates the country information.
func ValidateCountryInfo(ci *structs.CountryInfoInternal) error {
	err := validateCountryNameIsoCode(ci) // Validate the name and ISO code.
	if err != nil {
		return err
	}

	// Ensure that at least one feature is populated.
	if !ci.Features.Temperature && !ci.Features.Precipitation && !ci.Features.Capital &&
		!ci.Features.Coordinates && !ci.Features.Population && !ci.Features.Area &&
		(ci.Features.TargetCurrencies == nil || len(ci.Features.TargetCurrencies) == 0) {
		return errors.New("at least one feature must be populated")
	}
	return nil
}

// validateCountryNameIsoCode validates the provided country name and/or ISO code against supported countries.
func validateCountryNameIsoCode(ci *structs.CountryInfoInternal) error {
	validCountries, err := getSupportedCountries() // Fetch the list of supported countries.
	if err != nil {
		log.Printf("Error retriving supported countries: %v", err)
		return fmt.Errorf("error retriving supported countries: %v", err)
	}
	// Validate that a country has been specified.
	if err := validateCountryOrIsoCodeProvided(ci); err != nil {
		return err
	}
	// Validate the ISO, and update country name based on ISO. (If country name wasn't provided)
	if err := validateIsoCodeUpdateEmptyCountry(ci, validCountries); err != nil {
		return err
	}
	// Validate the country name, and update ISO based on country. (If ISO name wasn't provided)
	if err := updateIsoCodeAndValidateCountry(ci, validCountries); err != nil {
		return err
	}
	// Validate that the provided country and ISO correspond.
	return validateCorrespondence(ci, validCountries)
}

// validateCountryOrIsoCodeProvided checks that either country name or ISO code is provided.
func validateCountryOrIsoCodeProvided(ci *structs.CountryInfoInternal) error {
	if ci.Country == "" && ci.IsoCode == "" {
		return errors.New("either country name or ISO code must be provided")
	}
	return nil
}

// validateIsoCodeUpdateEmptyCountry checks if the provided ISO code is valid and updates the country name accordingly,
// if empty.
func validateIsoCodeUpdateEmptyCountry(ci *structs.CountryInfoInternal, validCountries map[string]string) error {
	if ci.IsoCode != "" { // Validate non-empty ISO.
		ci.IsoCode = strings.ToTitle(ci.IsoCode)                    // Convert ISO code to title case to match keys in validCountries.
		if country, exists := validCountries[ci.IsoCode]; !exists { // Validate ISO against validCountries.
			return errors.New("invalid ISO code")
		} else {
			if ci.Country == "" {
				ci.Country = country // Update the country name based on the ISO code, if empty.
			}
		}
	}
	return nil
}

// updateIsoCodeAndValidateCountry
// checks if only the country name was provided it validates it and updates the ISO code.
func updateIsoCodeAndValidateCountry(ci *structs.CountryInfoInternal, validCountries map[string]string) error {
	if ci.IsoCode == "" && ci.Country != "" {
		ci.Country = cases.Title(language.English, cases.Compact).String(ci.Country) // Normalize country name to title case.
		for code, name := range validCountries {
			if name == ci.Country {
				ci.IsoCode = code // Update ISO code if the country name is valid.
				return nil
			}
		}
		return errors.New("country name not valid or not supported")
	}
	return nil
}

// validateCorrespondence checks that the provided country name and ISO code match.
func validateCorrespondence(ci *structs.CountryInfoInternal, validCountries map[string]string) error {
	if ci.Country != "" && ci.IsoCode != "" {
		ci.Country = cases.Title(language.English, cases.Compact).String(ci.Country) // Normalize the country name.
		if validCountries[ci.IsoCode] != ci.Country {
			return errors.New("ISO code and country name do not match")
		}
	}
	return nil
}
