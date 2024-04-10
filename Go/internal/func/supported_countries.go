package _func

import (
	"encoding/json"
	"errors"
	"fmt"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"strings"
	"time"
)

// getSupportedCountries fetches countries with their common names and cca2 codes.
func getSupportedCountries() (map[string]string, error) {
	url := fmt.Sprintf("%sall?fields=name,cca2", External.CountriesAPI)
	var responseData []struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
		CCA2 string `json:"cca2"`
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error issuing request: %v", err)
	}

	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	countriesMap := make(map[string]string)
	for _, item := range responseData {
		countriesMap[item.CCA2] = item.Name.Common
	}
	return countriesMap, nil
}

func ValidateCountryInfo(ci *structs.CountryInfoGet) error {
	err := validateCountryNameIsoCode(ci)
	if err != nil {
		return err
	}

	if !ci.Features.Temperature && !ci.Features.Precipitation && !ci.Features.Capital &&
		!ci.Features.Coordinates && !ci.Features.Population && !ci.Features.Area {
		return errors.New("at least one feature must be true")
	}
	return nil
}

func validateCountryNameIsoCode(ci *structs.CountryInfoGet) error {
	validCountries, err := getSupportedCountries() // Adjusted to use the map version.
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

func validateCountryOrIsoCodeProvided(ci *structs.CountryInfoGet) error {
	if ci.Country == "" && ci.IsoCode == "" {
		return errors.New("either country name or ISO code must be provided")
	}
	return nil
}

func validateIsoCode(ci *structs.CountryInfoGet, validCountries map[string]string) error {
	if ci.IsoCode != "" {
		ci.IsoCode = strings.ToTitle(ci.IsoCode)
		if country, exists := validCountries[ci.IsoCode]; !exists {
			return errors.New("invalid ISO code")
		} else {
			ci.Country = country
		}
	}
	return nil
}

func updateAndValidateIsoCodeForCountry(ci *structs.CountryInfoGet, validCountries map[string]string) error {
	if ci.IsoCode == "" && ci.Country != "" {
		ci.Country = cases.Title(language.English, cases.Compact).String(ci.Country)
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

func validateCorrespondence(ci *structs.CountryInfoGet, validCountries map[string]string) error {
	if ci.Country != "" && ci.IsoCode != "" {
		ci.Country = cases.Title(language.English, cases.Compact).String(ci.Country)
		if validCountries[ci.IsoCode] != ci.Country {
			return errors.New("ISO code and country name do not match")
		}
	}
	return nil
}
