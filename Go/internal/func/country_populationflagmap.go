// Package _func provides self-defined functions for use in the API.
package _func

import (
	"encoding/json"
	"fmt"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"net/http"
)

// GetPopulationFlagMapForCountry retrieves population, maps, and flag information for a given country code.
// It returns the population count, Google Maps URL, PNG flag URL, and alternative text for the PNG flag.
func GetPopulationFlagMapForCountry(country string) (int, string, string, string) {
	// Create a new HTTP request to fetch population, maps, and flag data for the given country.
	r, err := http.NewRequest(http.MethodGet, External.CountriesAPI+"alpha?codes="+country+"&fields=population,maps,flags", nil)
	if err != nil {
		fmt.Printf("error creating request: %s", err.Error())
		return 0, "", "", ""
	}

	// Set content type header
	r.Header.Add("content-type", "application/json")

	// Create an HTTP client to execute the request.
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue the request and handle any errors.
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("error issuing request: %s", err.Error())
		return 0, "", "", ""
	}

	// Decode the JSON response into a list of Location structs.
	var locations []structs.Location
	err = json.NewDecoder(res.Body).Decode(&locations)
	if err != nil {
		fmt.Printf("error JSON decoding request: %s", err.Error())
		return 0, "", "", ""
	}

	// Ensure only one location is returned for the provided country.
	if len(locations) > 1 {
		fmt.Printf("Country parameter returned more than one package of population information, expected only one")
		return 0, "", "", ""
	}

	// Return population count, Google Maps URL, PNG flag URL, and alternative text for the PNG flag.
	return locations[0].Population, locations[0].Maps.GoogleMaps, locations[0].Flags.PNG, locations[0].Flags.ALT
}
