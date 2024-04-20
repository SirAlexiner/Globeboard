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

type OpenMeteoTemp struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

const (
	alphaCodes = "alpha?codes="
)

func GetTemp(coordinates structs.CoordinatesDashboard) (float64, error) {
	response, err := http.Get(External.OpenMeteoAPI + "?latitude=" + (coordinates.Latitude) + "&longitude=" + (coordinates.Longitude) + "&current=temperature_2m")
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing Open-Meteo api response body @ dashboardFunctions:GetTemp: ", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var openMeteo OpenMeteoTemp
	if err := json.Unmarshal(body, &openMeteo); err != nil {
		return 0, err
	}
	temp := openMeteo.Current.Temperature

	return temp, nil
}

type OpenMeteoPrecipitation struct {
	Current struct {
		Precipitation float64 `json:"precipitation"`
	} `json:"current"`
}

func GetPrecipitation(coordinates *structs.CoordinatesDashboard) (float64, error) {
	response, err := http.Get(External.OpenMeteoAPI + "?latitude=" + (coordinates.Latitude) + "&longitude=" + (coordinates.Longitude) + "&current=precipitation")
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing Open-Meteo api response body @ dashboardFunctions:GetTemp: ", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var openMeteo OpenMeteoPrecipitation
	if err := json.Unmarshal(body, &openMeteo); err != nil {
		return 0, err
	}

	precipitation := openMeteo.Current.Precipitation

	return precipitation, nil
}

type Country struct {
	Capital []string `json:"capital"`
}

func GetCapital(isocode string) (string, error) {
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=capital")
	if err != nil {
		return "nil", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing countries api response body @ dashboardFunctions:GetCapital: ", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var countries []Country
	if err := json.Unmarshal(body, &countries); err != nil {
		return "", err
	}
	capital := countries[0].Capital[0]

	return capital, nil
}

type CountryCoordinates struct {
	LatLng []float64 `json:"latlng"`
}

func GetCoordinates(isocode string) (structs.CoordinatesDashboard, error) {
	var empty = structs.CoordinatesDashboard{}

	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=latlng")
	if err != nil {
		return empty, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return empty, err
	}

	var countriesCoords []CountryCoordinates
	if err := json.Unmarshal(body, &countriesCoords); err != nil {
		return empty, err
	}

	var coords = structs.CoordinatesDashboard{
		Latitude:  strconv.FormatFloat(countriesCoords[0].LatLng[0], 'f', 5, 64),
		Longitude: strconv.FormatFloat(countriesCoords[0].LatLng[0], 'f', 5, 64),
	}

	return coords, nil
}

type CountryPopulation struct {
	Population int `json:"population"`
}

func GetPopulation(isocode string) (int, error) {
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=population")
	if err != nil {
		return 0, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var population []CountryPopulation
	if err := json.Unmarshal(body, &population); err != nil {
		return 0, err
	}
	populace := population[0].Population

	return populace, nil
}

type CountryArea struct {
	Area float64 `json:"area"`
}

func GetArea(isocode string) (float64, error) {
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=area")
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing countries api response body @ dashboardFunctions:GetArea: ", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var countryArea []CountryArea
	if err := json.Unmarshal(body, &countryArea); err != nil {
		return 0, err
	}
	area := countryArea[0].Area

	return area, nil
}

type CurrencyResponse []struct {
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

type RatesResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func GetExchangeRate(isocode string, currencies []string) (map[string]float64, error) {
	exchangeRateList, err := getExchangeRateList(isocode)
	if err != nil {
		return nil, err
	}
	exchangeRate := make(map[string]float64)
	for _, currency := range currencies {
		exchangeRate[strings.ToUpper(currency)] = exchangeRateList[strings.ToUpper(currency)]
	}

	return exchangeRate, nil
}

func fetchCurrencyRates(currency string) (map[string]float64, error) {
	response, err := http.Get(External.CurrencyAPI + currency)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing currency api response body: ", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var ratesData RatesResponse
	if err := json.Unmarshal(body, &ratesData); err != nil {
		return nil, err
	}

	return ratesData.Rates, nil
}

func getExchangeRateList(isocode string) (map[string]float64, error) {
	response, err := http.Get(External.CountriesAPI + alphaCodes + isocode + "&fields=currencies")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var currencyData CurrencyResponse
	if err := json.Unmarshal(body, &currencyData); err != nil {
		return nil, err
	}

	for currency := range currencyData[0].Currencies {
		rates, err := fetchCurrencyRates(currency)
		if err != nil {
			return nil, fmt.Errorf("error fetching currency rates: %v", err)
		}
		return rates, nil
	}

	return nil, errors.New("no currency data found")
}
