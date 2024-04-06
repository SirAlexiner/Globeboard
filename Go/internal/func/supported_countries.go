package _func

import (
	"encoding/json"
	"fmt"
	"globeboard/internal/utils/constants/External"
	"net/http"
	"time"
)

// CountryInfo represents the necessary information about a country, focusing on its common name and cca2 code.
type CountryInfo struct {
	CommonName string `json:"commonName"`
	ISOCode    string `json:"isoCode"`
}

// GetSupportedCountries fetches countries with their common names and cca2 codes.
func GetSupportedCountries() (map[string]string, error) {
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
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error issuing request: %s", err)
	}

	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %s", err)
	}

	countriesMap := make(map[string]string)
	for _, item := range responseData {
		countriesMap[item.CCA2] = item.Name.Common
	}
	return countriesMap, nil
}
