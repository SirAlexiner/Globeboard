// Package _func provides developer-made utility functions for use within the application.
package _func

import (
	"encoding/json"
	"errors"
	"fmt"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// OpenMeteoTemp structure defines the JSON structure for temperature response from the OpenMeteo API.
type OpenMeteoTemp struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"` // Current temperature (2 meters above the ground).
	} `json:"current"`
}

const (
	alphaCodes             = "alpha?codes="                    // URL parameter for filtering requests by country ISO codes.
	ResponseBodyCloseError = "Error closing response body: %v" // Log format for errors closing the response body.
)

// GetTemp fetches the current temperature for the specified coordinates using the OpenMeteo API.
func GetTemp(coordinates structs.CoordinatesDashboard) (float64, error) {
	// Constructing the URL to call the OpenMeteo API with query parameters for latitude and longitude.
	response, err := http.Get(External.OpenMeteoAPI + "?latitude=" + coordinates.Latitude + "&longitude=" + coordinates.Longitude + "&current=temperature_2m")
	if err != nil {
		log.Print(err)
		return 0, err // Return zero temperature and the error if the GET request fails.
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf(ResponseBodyCloseError, err)
		}
	}(response.Body) // Ensure the response body is closed after the function returns.

	body, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		log.Print(err)
		return 0, err // Return zero temperature and the error if an issue with reading the response occurs.
	}

	var openMeteo OpenMeteoTemp                              // Declaring a variable to store unmarshalled JSON data.
	if err := json.Unmarshal(body, &openMeteo); err != nil { // Unmarshal JSON data into OpenMeteoTemp struct.
		return 0, err
	}
	temp := openMeteo.Current.Temperature // Extracting the temperature from the struct.

	return temp, nil // Return the fetched temperature.
}

// OpenMeteoPrecipitation structure defines the JSON structure for precipitation response from OpenMeteo API.
type OpenMeteoPrecipitation struct {
	Current struct {
		Precipitation float64 `json:"precipitation"` // Current precipitation amount.
	} `json:"current"`
}

// GetPrecipitation fetches the current precipitation for the specified coordinates.
func GetPrecipitation(coordinates structs.CoordinatesDashboard) (float64, error) {
	// Construct the API request URL with coordinates.
	response, err := http.Get(External.OpenMeteoAPI + "?latitude=" + coordinates.Latitude + "&longitude=" + coordinates.Longitude + "&current=precipitation")
	if err != nil {
		log.Print(err)
		return 0, err // Return zero precipitations and the error if the GET request fails.
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf(ResponseBodyCloseError, err)
		}
	}(response.Body) // Ensuring the response body is properly closed.

	body, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		log.Print(err)
		return 0, err // Return zero precipitations and the error if an issue with reading the response occurs.
	}

	var openMeteo OpenMeteoPrecipitation                     // Struct to hold the precipitation data.
	if err := json.Unmarshal(body, &openMeteo); err != nil { // Parse the JSON response into the struct.
		return 0, err
	}

	precipitation := openMeteo.Current.Precipitation // Extract the precipitation value.

	return precipitation, nil // Return the precipitation data.
}

// Country structure for parsing country capital information from a JSON response.
type Country struct {
	Capital []string `json:"capital"` // Capital is expected as an array of strings,
	// though typically containing only one element.
}

// GetCapital fetches the capital city of a country identified by its ISO code.
func GetCapital(isocode string) (string, error) {
	// Construct the request URL with ISO code and fields parameter.
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=capital")
	if err != nil {
		log.Print(err)
		return "Earth", err // Return "Earth" and the error if GET requests fails.
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf(ResponseBodyCloseError, err)
		}
	}(response.Body) // Ensure to close the response body on function exit.

	body, err := io.ReadAll(response.Body) // Read the entire response body.
	if err != nil {
		log.Print(err)
		return "Earth", err //Return "Earth" and the error if an issue with reading the response occurs.
	}

	var countries []Country // Slice to hold the parsed JSON data.
	if err := json.Unmarshal(body, &countries); err != nil {
		return "", err // Handle JSON parsing errors.
	}
	if len(countries) == 0 || len(countries[0].Capital) == 0 {
		return "Earth", errors.New("no capital found for the specified ISO code") // Handle cases where no capital is available.
	}
	capital := countries[0].Capital[0] // Assume the first element is the desired capital.

	return capital, nil // Return the capital city.
}

// CountryCoordinates structure for parsing geographical coordinates from a JSON response.
type CountryCoordinates struct {
	LatLng []float64 `json:"latlng"` // Array of latitude and longitude.
}

// GetCoordinates fetches the geographical coordinates (latitude and longitude) of a country specified by its ISO code.
func GetCoordinates(isocode string) (structs.CoordinatesDashboard, error) {
	var empty = structs.CoordinatesDashboard{} // A default struct in case of errors.

	// Construct the request URL.
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=latlng")
	if err != nil {
		log.Print(err)
		return empty, err
	}

	body, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		log.Print(err)
		return empty, err
	}

	var countriesCoords []CountryCoordinates // Slice to store parsed data.
	if err := json.Unmarshal(body, &countriesCoords); err != nil {
		return empty, err // Handle JSON parsing errors.
	}

	// Assume the first entry contains the correct coordinates.
	var coords = structs.CoordinatesDashboard{
		Latitude:  strconv.FormatFloat(countriesCoords[0].LatLng[0], 'f', 5, 64),
		Longitude: strconv.FormatFloat(countriesCoords[0].LatLng[1], 'f', 5, 64),
	}

	return coords, nil // Return the coordinates in a struct.
}

// CountryPopulation structure for parsing population data from a JSON response.
type CountryPopulation struct {
	Population int `json:"population"` // Population as an integer.
}

// GetPopulation fetches the population of a country specified by its ISO code.
func GetPopulation(isocode string) (int, error) {
	// Construct the API request URL.
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=population")
	if err != nil {
		log.Print(err)
		return 0, err
	}

	body, err := io.ReadAll(response.Body) // Read the entire response body.
	if err != nil {
		log.Print(err)
		return 0, err
	}

	var populationData []CountryPopulation // Slice to hold parsed data.
	if err := json.Unmarshal(body, &populationData); err != nil {
		return 0, err // Handle JSON parsing errors.
	}

	if len(populationData) == 0 {
		return 0, errors.New("no population data found for the specified ISO code") // Handle cases where no data is found.
	}
	population := populationData[0].Population // Assume the first entry is the correct one.

	return population, nil // Return the population.
}

// CountryArea structure for parsing area data from a JSON response.
type CountryArea struct {
	Area float64 `json:"area"` // Area in square kilometers as a float64.
}

// GetArea fetches the total land area of a country specified by its ISO code.
func GetArea(isocode string) (float64, error) {
	// Construct the API request URL.
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=area")
	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf(ResponseBodyCloseError, err)
		}
	}(response.Body) // Ensure the response body is closed.

	body, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		log.Print(err)
		return 0, err
	}

	var countryArea []CountryArea // Slice to hold parsed data.
	if err := json.Unmarshal(body, &countryArea); err != nil {
		return 0, err // Handle JSON parsing errors.
	}

	if len(countryArea) == 0 {
		return 0, errors.New("no area data found for the specified ISO code") // Handle cases where no data is available.
	}
	area := countryArea[0].Area // Assume the first entry contains the correct area.

	return area, nil // Return the area.
}

// CurrencyResponse defines the structure for parsing currency information from a JSON response.
type CurrencyResponse []struct {
	Currencies map[string]struct {
		Name   string `json:"name"`   // Name of the currency.
		Symbol string `json:"symbol"` // Symbol of the currency.
	} `json:"currencies"`
}

// RatesResponse defines the structure for parsing exchange rate information from a JSON response.
type RatesResponse struct {
	Rates map[string]float64 `json:"rates"` // Map of currency codes to their respective exchange rates.
}

// GetExchangeRate computes the exchange rates for specified currencies against a base currency
// specified by the ISO code.
func GetExchangeRate(isocode string, currencies []string) (map[string]float64, error) {
	exchangeRateList, err := getExchangeRateList(isocode) // Fetch the list of all exchange rates for the base currency.
	if err != nil {
		log.Print(err)
		return nil, err
	}

	exchangeRate := make(map[string]float64) // Map to hold the filtered exchange rates.
	for _, currency := range currencies {
		exchangeRate[strings.ToUpper(currency)] = exchangeRateList[strings.ToUpper(currency)] // Filter and add the relevant rates.
	}

	return exchangeRate, nil // Return the map of exchange rates.
}

// fetchCurrencyRates retrieves the exchange rates for all currencies against a specified base currency.
func fetchCurrencyRates(currency string) (map[string]float64, error) {
	// Construct the API request URL.
	response, err := http.Get(External.CurrencyAPI + currency)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf(ResponseBodyCloseError, err)
		}
	}(response.Body) // Ensure the response body is closed.

	body, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var ratesData RatesResponse // Struct to hold the parsed data.
	if err := json.Unmarshal(body, &ratesData); err != nil {
		return nil, err // Handle JSON parsing errors.
	}

	return ratesData.Rates, nil // Return the map of exchange rates.
}

// getExchangeRateList fetches the exchange rates for all currencies against the base currency
// specified by the ISO code.
func getExchangeRateList(isocode string) (map[string]float64, error) {
	// Construct the API request URL for fetching currency information.
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=currencies")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	body, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var currencyData CurrencyResponse // Struct to hold the parsed currency information.
	if err := json.Unmarshal(body, &currencyData); err != nil {
		return nil, err // Handle JSON parsing errors.
	}

	for currency := range currencyData[0].Currencies {
		rates, err := fetchCurrencyRates(currency) // Fetch the exchange rates for the base currency.
		if err != nil {
			log.Printf("Error fetching currency rates: %v", err)
			return nil, fmt.Errorf("error fetching currency rates: %v", err) // Handle errors in fetching exchange rates.
		}
		return rates, nil // Return the map of exchange rates.
	}

	return nil, errors.New("no currency data found") // Return an error if no currency information is found.
}
